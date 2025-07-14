package api

import "net/http"

func PuushDelete(ctx *Context) {
	WritePuushError(ctx, -2, http.StatusNotImplemented)
}
