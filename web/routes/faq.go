package routes

import "github.com/Lekuruu/go-puush/internal/app"

func Faq(ctx *app.Context) {
	renderTemplate(ctx, "public/faq", map[string]interface{}{
		"Title": "faq",
	})
}
