package routes

import (
	"net/http"
	"strings"

	"github.com/Lekuruu/go-puush/internal/app"
)

func Login(ctx *app.Context) {
	renderTemplate(ctx, "public/login", map[string]interface{}{
		"Title": "login",
		"Retry": strings.Contains(ctx.Request.URL.Path, "retry"),
		"Error": ctx.Request.URL.Query().Get("error"),
	})
}

func PerformLogin(ctx *app.Context) {
	key := ctx.Request.URL.Query().Get("k")
	email := ctx.Request.FormValue("email")
	password := ctx.Request.FormValue("password")

	user, authenticated := UserAuthenticationDynamic(email, password, key, ctx.State)
	if !authenticated {
		http.Redirect(ctx.Response, ctx.Request, "/login/retry/", http.StatusSeeOther)
		return
	}

	if !user.Active {
		http.Redirect(ctx.Response, ctx.Request, "/login/retry/?error=inactive", http.StatusSeeOther)
		return
	}

	err := SetUserSession(user, ctx)
	if err != nil {
		http.Redirect(ctx.Response, ctx.Request, "/login/retry/?error=server", http.StatusSeeOther)
		return
	}

	http.Redirect(ctx.Response, ctx.Request, "/account", http.StatusFound)
}
