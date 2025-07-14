package api

import "net/http"

func PuushThumbnail(ctx *Context) {
	WritePuushError(ctx, -2, http.StatusNotImplemented)
}
