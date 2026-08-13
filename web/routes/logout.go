package routes

import (
	"net/http"

	"github.com/Lekuruu/go-puush/internal/server"
)

func Logout(ctx *server.Context) {
	// TODO: Ensure referer is safe before logging out
	//		 Else this could lead to CSRF logout "attacks"
	if err := ClearUserSession(ctx); err != nil {
		renderText(500, "failed to logout", ctx)
		return
	}

	http.Redirect(ctx.Response, ctx.Request, "/", http.StatusMovedPermanently)
}
