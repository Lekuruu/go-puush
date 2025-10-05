package routes

import "github.com/Lekuruu/go-puush/internal/app"

func About(ctx *app.Context) {
	renderTemplate(ctx, "public/about", map[string]interface{}{
		"Title": "about",
	})
}
