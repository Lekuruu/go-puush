package routes

import (
	"time"

	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/email"
	"github.com/Lekuruu/go-puush/internal/server"
	"github.com/Lekuruu/go-puush/internal/services"
	"github.com/Lekuruu/go-puush/internal/state"
)

func Register(ctx *server.Context) {
	renderTemplate(ctx, "public/register", map[string]any{
		"Title": "register",
	})
}

func PerformRegistration(ctx *server.Context) {
	if !ctx.State.Config.Service.RegistrationEnabled {
		renderErrorTemplate("Registrations Disabled", "Registrations are currently disabled. Please check back later!", ctx)
		return
	}

	err := ctx.Request.ParseForm()
	if err != nil {
		renderErrorTemplate("Error", "An error occurred while processing your registration.", ctx)
		return
	}

	email := ctx.Request.FormValue("email")
	password := ctx.Request.FormValue("password")
	confirm := ctx.Request.FormValue("confirm_password")

	if password != confirm {
		renderErrorTemplate("Password Mismatch", "The entered passwords do not match. Please try again!", ctx)
		return
	}

	if existingUser, _ := services.FetchUserByEmail(email, ctx.State); existingUser != nil {
		renderErrorTemplate("Email Taken", "Whoops! It appears that there is already an account with that email address.", ctx)
		return
	}

	if ctx.State.Config.Service.RequireInvitation {
		invitationKey := ctx.Request.FormValue("invitation_key")
		if invitationKey == "" {
			renderErrorTemplate("Missing Invitation Key", "An invitation key is required to register.", ctx)
			return
		}

		if valid, _ := services.IsValidInvitationKey(invitationKey, ctx.State); !valid {
			renderErrorTemplate("Invalid Invitation Key", "The provided invitation key is invalid or has expired.", ctx)
			return
		}
	}

	user, err := services.CreateUser(email, password, ctx.State)
	if err != nil {
		renderErrorTemplate("Error", "An error occurred while creating your account. Please try again!", ctx)
		return
	}

	ctx.Logger.Info("New user registered", "user_id", user.Id)

	// Write wal contents to disk after registration, if enabled
	if err := ctx.State.CheckpointWAL(); err != nil {
		ctx.Logger.Error("Failed to checkpoint database after registration", "error", err)
	}

	responseTitle := "Registration complete!"
	responseMessage := "You can now log in with your email and password."

	if ctx.State.Config.Service.RequireActivation {
		responseMessage = "An email has been sent to your designated address with instructions on how to activate your account."
		go createAndSendActivationEmail(user, ctx.State)
	}

	renderResponseTemplate(responseTitle, responseMessage, "registration complete", ctx)
}

func PerformActivation(ctx *server.Context) {
	if !ctx.State.Config.Service.RegistrationEnabled {
		renderErrorTemplate("Sorry! Account activation is currently disabled.", "Please contact support for assistance.", ctx)
		return
	}

	key := ctx.Request.URL.Query().Get("key")
	if key == "" {
		renderErrorTemplate("Uh-oh! Looks like you missed a chunk of your verification link!", "Please make sure you've copied the whole URL from your email and try again.", ctx)
		return
	}

	verification, err := services.FetchEmailVerificationByKey(key, ctx.State, "User")
	if err != nil || verification.Action != database.EmailVerificationActionActivate {
		renderErrorTemplate("Uh-oh! Looks like this verification link is invalid.", "Please make sure you've entered the URL correctly and try again.", ctx)
		return
	}

	if verification.User.Active {
		renderErrorTemplate("This account is already verified.", "Please log in to your account.", ctx)
		return
	}

	verification.User.Active = true
	err = services.ActivateUser(verification.User.Id, ctx.State)
	if err != nil {
		renderErrorTemplate("Uh-oh! Something went wrong.", "An error occurred while activating your account. Please try again later.", ctx)
		return
	}

	go sendWelcomeEmail(verification.User, ctx.State)
	services.DeleteEmailVerificationById(verification.Id, ctx.State)
	renderResponseTemplate("Activation complete!", "Your account has been successfully activated. You can now log in.", "activation complete", ctx)
}

const emailVerificationExpiry = time.Hour * 24 * 7

func createAndSendActivationEmail(user *database.User, state *state.State) {
	verification, err := services.CreateEmailVerification(&user.Id, database.EmailVerificationActionActivate, emailVerificationExpiry, state)
	if err != nil {
		state.Logger.Error("Failed to create email verification", "user_id", user.Id, "error", err)
		return
	}

	message := email.FormatActivationEmail(user.Email, verification.Key, state.Config.Service.Url)
	err = state.Email.Send(message)
	if err != nil {
		state.Logger.Error("Failed to send account activation email", "user_id", user.Id, "error", err)
		return
	}
}

func sendWelcomeEmail(user *database.User, state *state.State) {
	message := email.FormatWelcomeEmail(user.Email)
	err := state.Email.Send(message)
	if err != nil {
		state.Logger.Error("Failed to send welcome email", "user_id", user.Id, "error", err)
	}
}
