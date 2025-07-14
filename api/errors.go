package api

import "net/http"

func PuushErrorSubmission(ctx *Context) {
	WritePuushError(ctx, -2, http.StatusNotImplemented)
}
