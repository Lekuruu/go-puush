package routes

import "github.com/Lekuruu/go-puush/internal/server"

func Dmca(ctx *server.Context) {
	renderTemplate(ctx, "public/dmca", map[string]any{
		"Title": "DMCA",
	})
}
