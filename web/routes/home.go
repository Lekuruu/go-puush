package routes

import "github.com/Lekuruu/go-puush/internal/app"

func Home(ctx *app.Context) {
	renderTemplate(ctx, "public/home", map[string]interface{}{
		"Title": "home",
	})
}
