package main

import (
	"github.com/Lekuruu/go-puush/api"
	"github.com/Lekuruu/go-puush/internal/app"
)

func InitializeRoutes(server *app.Server) {
	server.Router.HandleFunc("/api/register", server.ContextMiddleware(api.PuushRegistration)).Methods("POST")
	server.Router.HandleFunc("/api/auth", server.ContextMiddleware(api.PuushAuthentication)).Methods("POST")
	server.Router.HandleFunc("/api/up", server.ContextMiddleware(api.PuushUpload)).Methods("POST")
	server.Router.HandleFunc("/api/del", server.ContextMiddleware(api.PuushDelete)).Methods("POST")
	server.Router.HandleFunc("/api/hist", server.ContextMiddleware(api.PuushHistory)).Methods("POST")
	server.Router.HandleFunc("/api/thumb", server.ContextMiddleware(api.PuushThumbnail)).Methods("POST")
	server.Router.HandleFunc("/api/oshi", server.ContextMiddleware(api.PuushErrorSubmission)).Methods("POST")
	server.Router.HandleFunc("/dl/puush-rss.xml", server.ContextMiddleware(api.PuushMacOSRssFeed)).Methods("GET")
	server.Router.HandleFunc("/dl/puush-win.txt", server.ContextMiddleware(api.PuushWindowsUpdate)).Methods("GET")
}

func main() {
	state := app.NewState()
	if state == nil {
		return
	}

	server := app.NewServer(
		state.Config.Api.Host,
		state.Config.Api.Port,
		"puush-api",
		state,
	)
	InitializeRoutes(server)
	server.Serve()
}
