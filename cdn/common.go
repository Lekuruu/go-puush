package cdn

import "github.com/Lekuruu/go-puush/internal/app"

func WriteResponse(status int, message string, ctx *app.Context) {
	ctx.Response.WriteHeader(status)
	ctx.Response.Header().Set("Content-Type", "text/plain")
	ctx.Response.Write([]byte(message))
}
