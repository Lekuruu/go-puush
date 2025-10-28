package routes

import (
	"time"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/email"
	"github.com/Lekuruu/go-puush/internal/services"
)

func Register(ctx *app.Context) {
	renderTemplate(ctx, "public/register", map[string]interface{}{
		"Title": "register",
	})
}

func PerformRegistration(ctx *app.Context) {
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

	ctx.State.Logger.Logf(
		"New user registered: %s (%d)",
		user.Email, user.Id,
	)

	responseTitle := "Registration complete!"
	responseMessage := "You can now log in with your email and password."

	if ctx.State.Config.Service.RequireActivation {
		responseMessage = "An email has been sent to your designated address with instructions on how to activate your account."
		go createAndSendActivationEmail(user, ctx.State)
	}

	renderResponseTemplate(responseTitle, responseMessage, "registration complete", ctx)
}

func PerformActivation(ctx *app.Context) {
	if !ctx.State.Config.Service.RegistrationEnabled {
		renderErrorTemplate("Activation Disabled", "Account activation is currently disabled. Please contact support for assistance.", ctx)
		return
	}

	key := ctx.Request.URL.Query().Get("key")
	if key == "" {
		renderErrorTemplate("Invalid Activation Key", "The activation key provided is invalid. Please check your email for the correct link.", ctx)
		return
	}

	verification, err := services.FetchEmailVerificationByKey(key, ctx.State, "User")
	if err != nil || verification.Action != database.EmailVerificationActionActivate {
		renderErrorTemplate("Invalid Activation Key", "The activation key provided is invalid or has already been used.", ctx)
		return
	}

	if verification.User.Active {
		renderErrorTemplate("Already Activated", "Your account is already activated. You can log in directly.", ctx)
		return
	}

	verification.User.Active = true
	err = services.UpdateUser(verification.User, ctx.State)
	if err != nil {
		renderErrorTemplate("Activation Error", "An error occurred while activating your account. Please try again later.", ctx)
		return
	}

	services.DeleteEmailVerificationById(verification.Id, ctx.State)
	renderResponseTemplate("Activation complete!", "Your account has been successfully activated. You can now log in.", "activation complete", ctx)
}

const emailVerificationExpiry = time.Hour * 24 * 7

func createAndSendActivationEmail(user *database.User, state *app.State) {
	verification, err := services.CreateEmailVerification(&user.Id, database.EmailVerificationActionActivate, emailVerificationExpiry, state)
	if err != nil {
		state.Logger.Logf("Failed to create email verification for user %d: %v", user.Id, err)
		return
	}

	message := email.FormatActivationEmail(user.Email, verification.Key)
	err = state.Email.Send(message)
	if err != nil {
		state.Logger.Logf("Failed to send account activation email to user %d: %v", user.Id, err)
		return
	}
}
