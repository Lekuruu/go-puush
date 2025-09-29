package routes

import "github.com/Lekuruu/go-puush/internal/app"

func Login(ctx *app.Context) {
	renderTemplate(ctx, "public/login", map[string]interface{}{
		"Title": "login",
	})
}
