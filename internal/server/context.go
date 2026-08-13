package server

import (
	"log/slog"
	"net/http"

	"github.com/Lekuruu/go-puush/internal/state"
)

// Context is a struct that holds the request context for each endpoint call.
type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	State    *state.State
	Logger   *slog.Logger
}

func (ctx *Context) IP() string {
	return GetRequestIP(ctx.Request)
}

func (ctx *Context) PathValue(name string) string {
	return ctx.Request.PathValue(name)
}
