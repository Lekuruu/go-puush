package api

import "github.com/Lekuruu/go-puush/internal/app"

func PuushErrorSubmission(ctx *app.Context) {
	WritePuushError(ctx, NotImplementedError)
}
