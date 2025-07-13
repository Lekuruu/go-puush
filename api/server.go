package api

import (
	"fmt"
	"net/http"

	"github.com/Lekuruu/go-puush/api/puush"
	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/gorilla/mux"
)

type Server struct {
	State  *app.State
	Router *mux.Router
}

func NewServer(state *app.State) *Server {
	return &Server{State: state}
}

func (server *Server) Serve() {
	bind := fmt.Sprintf(
		"%s:%d",
		server.State.Config.Server.Host,
		server.State.Config.Server.Port,
	)

	server.Router = mux.NewRouter()
	server.Router.HandleFunc("/api/auth", server.ContextMiddleware(puush.PuushAuthentication)).Methods("POST")
	server.Router.HandleFunc("/api/up", server.ContextMiddleware(puush.PuushUpload)).Methods("POST")
	server.Router.HandleFunc("/api/del", server.ContextMiddleware(puush.PuushDelete)).Methods("POST")
	server.Router.HandleFunc("/api/hist", server.ContextMiddleware(puush.PuushHistory)).Methods("POST")
	server.Router.HandleFunc("/api/thumb", server.ContextMiddleware(puush.PuushThumbnail)).Methods("POST")
	server.Router.HandleFunc("/api/oshi", server.ContextMiddleware(puush.PuushErrorSubmission)).Methods("POST")
	http.ListenAndServe(bind, server.Router)
}
