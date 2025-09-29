package routes

import "github.com/Lekuruu/go-puush/internal/app"

func Register(ctx *app.Context) {
	renderTemplate(ctx, "public/register", map[string]interface{}{
		"Title": "register",
	})
}
