package routes

import "github.com/Lekuruu/go-puush/internal/app"

func TermsOfService(ctx *app.Context) {
	renderTemplate(ctx, "public/tos", map[string]interface{}{
		"Title": "Terms of Service",
	})
}
