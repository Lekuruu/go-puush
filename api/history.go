package api

import "net/http"

func PuushHistory(ctx *Context) {
	user, err := UserAuthenticationFromContext(ctx)
	if err != nil {
		WritePuushError(ctx, -1, http.StatusUnauthorized)
		return
	}

	if user == nil {
		WritePuushError(ctx, -2, http.StatusInternalServerError)
		return
	}

	// TODO: Retrieve uploads
}
