package api

import "net/http"

func PuushRegistration(ctx *Context) {
	WritePuushError(ctx, -2, http.StatusNotImplemented)
}
