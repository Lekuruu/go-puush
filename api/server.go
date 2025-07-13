package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/gorilla/mux"
)

type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	State    *app.State
}

type Server struct {
	State  *app.State
	Router *mux.Router
	Logger *slog.Logger
}

func NewServer(state *app.State) *Server {
	return &Server{
		State:  state,
		Logger: state.Logger,
	}
}

func (server *Server) Serve() {
	server.Router = mux.NewRouter()
	server.InitializeRoutes()

	bind := fmt.Sprintf(
		"%s:%d",
		server.State.Config.Api.Host,
		server.State.Config.Api.Port,
	)
	server.Logger.Info(
		"Starting server",
		"bind", bind,
	)

	err := http.ListenAndServe(bind, server.Router)
	if err != nil {
		slog.Error("Failed to start server", "error", err)
		return
	}
}

func (server *Server) InitializeRoutes() {
	server.Router.HandleFunc("/api/auth", server.ContextMiddleware(PuushAuthentication)).Methods("POST")
	server.Router.HandleFunc("/api/up", server.ContextMiddleware(PuushUpload)).Methods("POST")
	server.Router.HandleFunc("/api/del", server.ContextMiddleware(PuushDelete)).Methods("POST")
	server.Router.HandleFunc("/api/hist", server.ContextMiddleware(PuushHistory)).Methods("POST")
	server.Router.HandleFunc("/api/thumb", server.ContextMiddleware(PuushThumbnail)).Methods("POST")
	server.Router.HandleFunc("/api/oshi", server.ContextMiddleware(PuushErrorSubmission)).Methods("POST")
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
