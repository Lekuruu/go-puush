package main

import (
	"net/http"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/web/routes"
)

func InitializeRoutes(server *app.Server) {
	// Public pages
	server.Router.HandleFunc("/", server.ContextMiddleware(routes.Home)).Methods("GET")
	server.Router.HandleFunc("/faq", server.ContextMiddleware(routes.Faq)).Methods("GET")
	server.Router.HandleFunc("/about", server.ContextMiddleware(routes.About)).Methods("GET")
	server.Router.HandleFunc("/register", server.ContextMiddleware(routes.Register)).Methods("GET")
	server.Router.HandleFunc("/reset_password", server.ContextMiddleware(routes.ResetPassword)).Methods("GET")
	server.Router.HandleFunc("/tos", server.ContextMiddleware(routes.TermsOfService)).Methods("GET")
	server.Router.HandleFunc("/dmca", server.ContextMiddleware(routes.Dmca)).Methods("GET")

	// Login pages
	server.Router.HandleFunc("/login", server.ContextMiddleware(routes.Login)).Methods("GET")
	server.Router.HandleFunc("/login/", server.ContextMiddleware(routes.Login)).Methods("GET")
	server.Router.HandleFunc("/login/go", server.ContextMiddleware(routes.PerformLogin)).Methods("POST")
	server.Router.HandleFunc("/login/retry", server.ContextMiddleware(routes.Login)).Methods("GET")
	server.Router.HandleFunc("/login/retry/", server.ContextMiddleware(routes.Login)).Methods("GET")
	server.Router.HandleFunc("/logout", server.ContextMiddleware(routes.Logout)).Methods("GET")
	server.Router.HandleFunc("/logout/", server.ContextMiddleware(routes.Logout)).Methods("GET")

	// Gallery pages
	server.Router.HandleFunc("/{username}/Gallery", server.ContextMiddleware(routes.Gallery)).Methods("GET")
	server.Router.HandleFunc("/{username}/Public/feed", server.ContextMiddleware(routes.GalleryFeed)).Methods("GET")

	// Account pages
	server.Router.HandleFunc("/account", server.ContextMiddleware(routes.Account)).Methods("GET")
	server.Router.HandleFunc("/account/", server.ContextMiddleware(routes.Account)).Methods("GET")
	server.Router.HandleFunc("/account/go_pro", server.ContextMiddleware(routes.AccountGoPro)).Methods("GET")
	server.Router.HandleFunc("/account/settings", server.ContextMiddleware(routes.AccountSettings)).Methods("GET")
	server.Router.HandleFunc("/account/subscription", server.ContextMiddleware(routes.AccountSubscription)).Methods("GET")

	// AJAX pages
	server.Router.HandleFunc("/ajax/move_dialog/", server.ContextMiddleware(routes.MoveDialog)).Methods("GET")
	server.Router.HandleFunc("/ajax/move_upload", server.ContextMiddleware(routes.MoveUpload)).Methods("POST")

	// Thumbnail page
	server.Router.HandleFunc("/thumb/view/{identifier}", server.ContextMiddleware(routes.Thumbnail)).Methods("GET")

	// Initialize static routes
	server.Router.PathPrefix("/dl/").Handler(http.StripPrefix("/dl/", http.FileServer(http.Dir("web/static/dl/"))))
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
		state.Config.Web.Host,
		state.Config.Web.Port,
		"puush-web",
		state,
	)
	InitializeRoutes(server)
	server.Serve()
}
