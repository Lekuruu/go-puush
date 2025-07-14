package api

import "net/http"

func PuushUpload(ctx *Context) {
	WritePuushError(ctx, -2, http.StatusNotImplemented)
}
