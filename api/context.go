package api

import (
	"net/http"

	"github.com/Lekuruu/go-puush/internal/app"
)

type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	State    *app.State
}

func (server *Server) ContextMiddleware(handler func(*Context)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		context := &Context{
			Response: w,
			Request:  r,
			State:    server.State,
		}

		w.Header().Set("Server", "puush")
		handler(context)
	}
}
