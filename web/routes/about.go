package routes

import "github.com/Lekuruu/go-puush/internal/server"

func About(ctx *server.Context) {
	renderTemplate(ctx, "public/about", map[string]any{
		"Title": "about",
	})
}
