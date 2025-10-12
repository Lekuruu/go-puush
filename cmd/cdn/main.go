package main

import (
	"github.com/Lekuruu/go-puush/cdn"
	"github.com/Lekuruu/go-puush/internal/app"
)

func InitializeRoutes(server *app.Server) {
	server.Router.HandleFunc("/{identifier}", server.ContextMiddleware(cdn.UploadShortlink)).Methods("GET")
	server.Router.HandleFunc("/t/{identifier}", server.ContextMiddleware(cdn.ThumbnailShortlink)).Methods("GET")
	server.Router.HandleFunc("/{pool}/{identifier}", server.ContextMiddleware(cdn.UploadShortlink)).Methods("GET")
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
