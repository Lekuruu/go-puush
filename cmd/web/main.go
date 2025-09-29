package main

import (
	"net/http"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/web/routes"
)

func InitializeRoutes(server *app.Server) {
	server.Router.HandleFunc("/", server.ContextMiddleware(routes.Home)).Methods("GET")
	server.Router.HandleFunc("/about", server.ContextMiddleware(routes.About)).Methods("GET")
	server.Router.HandleFunc("/login", server.ContextMiddleware(routes.Login)).Methods("GET")
	server.Router.HandleFunc("/account", server.ContextMiddleware(routes.Account)).Methods("GET")
	server.Router.HandleFunc("/register", server.ContextMiddleware(routes.Register)).Methods("GET")
	server.Router.HandleFunc("/reset_password", server.ContextMiddleware(routes.ResetPassword)).Methods("GET")
	server.Router.HandleFunc("/tos", server.ContextMiddleware(routes.TermsOfService)).Methods("GET")
	server.Router.HandleFunc("/dmca", server.ContextMiddleware(routes.Dmca)).Methods("GET")

	// Initialize static routes
	server.Router.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("web/static/js/"))))
	server.Router.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("web/static/css/"))))
	server.Router.PathPrefix("/img/").Handler(http.StripPrefix("/img/", http.FileServer(http.Dir("web/static/img/"))))
	server.Router.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/static/robots.txt")
	})
	server.Router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/static/favicon.ico")
	})
}

func main() {
	state := app.NewState()
	if state == nil {
		return
	}

	server := app.NewServer(
		state.Config.Cdn.Host,
		state.Config.Cdn.Port,
		"puush-web",
		state,
	)
	InitializeRoutes(server)
	server.Serve()
}
