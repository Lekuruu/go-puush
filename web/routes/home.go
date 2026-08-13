package routes

import "github.com/Lekuruu/go-puush/internal/server"

func Home(ctx *server.Context) {
	renderTemplate(ctx, "public/home", map[string]any{
		"Title": "home",
	})
}
