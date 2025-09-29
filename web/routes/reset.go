package routes

import "github.com/Lekuruu/go-puush/internal/app"

func ResetPassword(ctx *app.Context) {
	renderTemplate(ctx, "public/reset", map[string]interface{}{
		"Title": "reset password",
	})
}
