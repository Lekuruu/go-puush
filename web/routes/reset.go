package routes

import (
	"net/http"
	"strings"
	"time"

	"github.com/Lekuruu/go-puush/internal/authentication"
	"github.com/Lekuruu/go-puush/internal/database"
	"github.com/Lekuruu/go-puush/internal/email"
	"github.com/Lekuruu/go-puush/internal/server"
	"github.com/Lekuruu/go-puush/internal/services"
	"github.com/Lekuruu/go-puush/internal/state"
)

func ResetPassword(ctx *server.Context) {
	renderTemplate(ctx, "public/reset", map[string]any{
		"Title": "reset password",
	})
}

func RequestPasswordReset(ctx *server.Context) {
	if err := ctx.Request.ParseForm(); err != nil {
		renderErrorTemplate("Uh-oh! Something went wrong.", "An error occurred while submitting your request. Please try again!", ctx)
		return
	}

	emailAddress := strings.TrimSpace(ctx.Request.FormValue("email"))
	if emailAddress == "" {
		renderErrorTemplate("Missing Email", "Please enter the email address associated with your account and try again!", ctx)
		return
	}

	user, err := services.FetchUserByEmail(emailAddress, ctx.State)
	if err != nil {
		renderErrorTemplate("Invalid Account", "Sorry, that account was either not found or has not yet been activated. Please recheck it and try again.", ctx)
		return
	}

	if !user.Active {
		renderErrorTemplate("Invalid Account", "Sorry, that account was either not found or has not yet been activated. Please recheck it and try again.", ctx)
		return
	}

	go createAndSendPasswordResetEmail(user, ctx.State)

	// Write wal contents to disk, if enabled
	if err := ctx.State.CheckpointWAL(); err != nil {
		ctx.Logger.Error("Failed to checkpoint database after password reset request", "error", err)
	}

	renderResponseTemplate(
		"Password reset request received!",
		"We'll send an email with your password reset link shortly.",
		"password reset",
		ctx,
	)
}

func ShowPasswordResetForm(ctx *server.Context) {
	key := strings.TrimSpace(ctx.Request.URL.Query().Get("key"))
	if key == "" {
		http.Redirect(ctx.Response, ctx.Request, "/reset_password", http.StatusSeeOther)
		return
	}

	verification, err := services.ValidateEmailVerification(
		key,
		database.EmailVerificationActionResetPassword,
		ctx.State,
	)
	if err != nil {
		ctx.Logger.Error("Failed to validate password reset key", "error", err)
		renderErrorTemplate("Uh-oh! Something went wrong.", "Please try requesting a new password reset link.", ctx)
		return
	}

	if verification == nil {
		http.Redirect(ctx.Response, ctx.Request, "/reset_password", http.StatusSeeOther)
		return
	}

	renderTemplate(ctx, "public/reset/password_select", map[string]any{
		"Title": "reset password",
		"Key":   verification.Key,
	})
}

func PerformPasswordReset(ctx *server.Context) {
	if err := ctx.Request.ParseForm(); err != nil {
		renderErrorTemplate("Uh-oh! Something went wrong.", "An error occurred while submitting your request. Please try again!", ctx)
		return
	}

	key := strings.TrimSpace(ctx.Request.FormValue("key"))
	password := ctx.Request.FormValue("password")
	confirm := ctx.Request.FormValue("confirm_password")

	if key == "" {
		renderErrorTemplate("Missing Reset Key", "It looks like your reset link is incomplete. Please use the link from your email and try again.", ctx)
		return
	}

	if password == "" || confirm == "" {
		renderErrorTemplate("Missing Password", "Please enter your new password twice to confirm the change.", ctx)
		return
	}

	if password != confirm {
		renderErrorTemplate("Password Mismatch", "The entered passwords do not match. Please try again!", ctx)
		return
	}

	verification, err := services.ValidateEmailVerification(
		key,
		database.EmailVerificationActionResetPassword,
		ctx.State,
		"User",
	)
	if err != nil {
		ctx.Logger.Error("Failed to validate password reset key", "error", err)
		renderErrorTemplate("Uh-oh! Something went wrong.", "Please try requesting a new password reset link.", ctx)
		return
	}

	newPasswordHash, err := authentication.CreatePasswordHash(password)
	if err != nil {
		ctx.Logger.Error("Failed to hash new password", "user_id", verification.User.Id, "error", err)
		renderErrorTemplate("Uh-oh! Something went wrong.", "Please try again later.", ctx)
		return
	}

	if err := services.UpdateUserPassword(verification.User.Id, newPasswordHash, ctx.State); err != nil {
		ctx.Logger.Error("Failed to update password", "user_id", verification.User.Id, "error", err)
		renderErrorTemplate("Uh-oh! Something went wrong.", "Please try again later.", ctx)
		return
	}

	if err := services.DeleteEmailVerificationById(verification.Id, ctx.State); err != nil {
		ctx.Logger.Error("Failed to delete password reset verification", "verification_id", verification.Id, "error", err)
	}

	// Write wal contents to disk, if enabled
	if err := ctx.State.CheckpointWAL(); err != nil {
		ctx.Logger.Error("Failed to checkpoint database after password reset", "error", err)
	}

	renderResponseTemplate(
		"Password reset complete!",
		"Your password has been updated. You can now log in with your new credentials.",
		"password reset",
		ctx,
	)
}

const passwordResetVerificationExpiry = time.Hour

func createAndSendPasswordResetEmail(user *database.User, state *state.State) {
	verification, err := services.CreateEmailVerification(&user.Id, database.EmailVerificationActionResetPassword, passwordResetVerificationExpiry, state)
	if err != nil {
		state.Logger.Error("Failed to create password reset verification", "user_id", user.Id, "error", err)
		return
	}

	message := email.FormatPasswordResetEmail(user.Email, verification.Key, state.Config.Service.Url)
	if err := state.Email.Send(message); err != nil {
		state.Logger.Error("Failed to send password reset email", "user_id", user.Id, "error", err)
	}
}
