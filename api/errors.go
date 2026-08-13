package api

import "github.com/Lekuruu/go-puush/internal/server"

func PuushErrorSubmission(ctx *server.Context) {
	WritePuushError(ctx, NotImplementedError)
}
