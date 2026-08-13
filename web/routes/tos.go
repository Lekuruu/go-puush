package routes

import "github.com/Lekuruu/go-puush/internal/server"

func TermsOfService(ctx *server.Context) {
	renderTemplate(ctx, "public/tos", map[string]any{
		"Title": "Terms of Service",
	})
}
