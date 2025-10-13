package routes

import (
	"github.com/Lekuruu/go-puush/internal/app"
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
		renderErrorTemplate("Email Taken", "An account with that email address already exists.", ctx)
		return
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
	responseMessage := "An email has been sent to your designated address with instructions on how to activate your account."

	if !ctx.State.Config.Service.RequireActivation {
		responseMessage = "You can now log in with your email and password."
	}

	// TODO: Send activation email
	renderResponseTemplate(responseTitle, responseMessage, "registration complete", ctx)
}
