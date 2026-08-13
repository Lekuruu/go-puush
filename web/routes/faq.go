package routes

import "github.com/Lekuruu/go-puush/internal/server"

func Faq(ctx *server.Context) {
	renderTemplate(ctx, "public/faq", map[string]any{
		"Title": "faq",
	})
}
