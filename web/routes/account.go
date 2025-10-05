package routes

import (
	"net/http"

	"github.com/Lekuruu/go-puush/internal/app"
)

func Account(ctx *app.Context) {
	// Redirect account page to login for now
	http.Redirect(ctx.Response, ctx.Request, "/login", http.StatusTemporaryRedirect)
}
