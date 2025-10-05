package routes

import "github.com/Lekuruu/go-puush/internal/app"

func Dmca(ctx *app.Context) {
	renderTemplate(ctx, "public/dmca", map[string]interface{}{
		"Title": "DMCA",
	})
}
