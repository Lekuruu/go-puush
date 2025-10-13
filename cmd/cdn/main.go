package main

import (
	"github.com/Lekuruu/go-puush/cdn"
	"github.com/Lekuruu/go-puush/internal/app"
)

func InitializeRoutes(server *app.Server) {
	server.Router.HandleFunc("/{identifier}", server.ContextMiddleware(cdn.Upload)).Methods("GET")
	server.Router.HandleFunc("/{pool}/{identifier}", server.ContextMiddleware(cdn.Upload)).Methods("GET")
	server.Router.HandleFunc("/t/{identifier}", server.ContextMiddleware(cdn.Thumbnail)).Methods("GET")
	server.Router.HandleFunc("/t/{pool}/{identifier}", server.ContextMiddleware(cdn.Thumbnail)).Methods("GET")
}

func main() {
	state := app.NewState()
	if state == nil {
		return
	}

	server := app.NewServer(
		state.Config.Cdn.Host,
		state.Config.Cdn.Port,
		"puush-cdn",
		state,
	)
	InitializeRoutes(server)
	server.Serve()
}
