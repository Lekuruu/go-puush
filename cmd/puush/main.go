package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Lekuruu/go-puush/api"
	"github.com/Lekuruu/go-puush/cdn"
	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/Lekuruu/go-puush/web/routes"
)

func InitializeCdnRoutes(server *app.Server) {
	server.Router.HandleFunc("/{identifier}", server.ContextMiddleware(cdn.Upload)).Methods("GET")
	server.Router.HandleFunc("/{pool}/{identifier}", server.ContextMiddleware(cdn.Upload)).Methods("GET")
	server.Router.HandleFunc("/t/{identifier}", server.ContextMiddleware(cdn.Thumbnail)).Methods("GET")
	server.Router.HandleFunc("/t/{pool}/{identifier}", server.ContextMiddleware(cdn.Thumbnail)).Methods("GET")
}

func InitializeApiRoutes(server *app.Server) {
	server.Router.HandleFunc("/api/register", server.ContextMiddleware(api.PuushRegistration)).Methods("POST")
	server.Router.HandleFunc("/api/auth", server.ContextMiddleware(api.PuushAuthentication)).Methods("POST")
	server.Router.HandleFunc("/api/up", server.ContextMiddleware(api.PuushUpload)).Methods("POST")
	server.Router.HandleFunc("/api/del", server.ContextMiddleware(api.PuushDelete)).Methods("POST")
	server.Router.HandleFunc("/api/hist", server.ContextMiddleware(api.PuushHistory)).Methods("POST")
	server.Router.HandleFunc("/api/thumb", server.ContextMiddleware(api.PuushThumbnail)).Methods("POST")
	server.Router.HandleFunc("/api/oshi", server.ContextMiddleware(api.PuushErrorSubmission)).Methods("POST")
}

func InitializeWebRoutes(server *app.Server) {
	// Initialize templates
	routes.InitializeTemplates()

	// Public pages
	server.Router.HandleFunc("/", server.ContextMiddleware(routes.Home)).Methods("GET")
	server.Router.HandleFunc("/faq", server.ContextMiddleware(routes.Faq)).Methods("GET")
	server.Router.HandleFunc("/about", server.ContextMiddleware(routes.About)).Methods("GET")
	server.Router.HandleFunc("/register", server.ContextMiddleware(routes.Register)).Methods("GET")
	server.Router.HandleFunc("/register/go", server.ContextMiddleware(routes.PerformRegistration)).Methods("POST")
	server.Router.HandleFunc("/register/verify", server.ContextMiddleware(routes.PerformActivation)).Methods("GET")
	server.Router.HandleFunc("/reset_password", server.ContextMiddleware(routes.ResetPassword)).Methods("GET")
	server.Router.HandleFunc("/reset_password", server.ContextMiddleware(routes.RequestPasswordReset)).Methods("POST")
	server.Router.HandleFunc("/reset_password/go", server.ContextMiddleware(routes.ShowPasswordResetForm)).Methods("GET")
	server.Router.HandleFunc("/reset_password/go", server.ContextMiddleware(routes.PerformPasswordReset)).Methods("POST")
	server.Router.HandleFunc("/tos", server.ContextMiddleware(routes.TermsOfService)).Methods("GET")
	server.Router.HandleFunc("/dmca", server.ContextMiddleware(routes.Dmca)).Methods("GET")

	// Login pages
	server.Router.HandleFunc("/login", server.ContextMiddleware(routes.Login)).Methods("GET")
	server.Router.HandleFunc("/login/", server.ContextMiddleware(routes.Login)).Methods("GET")
	server.Router.HandleFunc("/login/go", server.ContextMiddleware(routes.PerformLogin)).Methods("GET", "POST")
	server.Router.HandleFunc("/login/go/", server.ContextMiddleware(routes.PerformLogin)).Methods("GET", "POST")
	server.Router.HandleFunc("/login/retry", server.ContextMiddleware(routes.Login)).Methods("GET")
	server.Router.HandleFunc("/login/retry/", server.ContextMiddleware(routes.Login)).Methods("GET")
	server.Router.HandleFunc("/logout", server.ContextMiddleware(routes.Logout)).Methods("GET")
	server.Router.HandleFunc("/logout/", server.ContextMiddleware(routes.Logout)).Methods("GET")

	// Gallery pages
	server.Router.HandleFunc("/{username}/Gallery", server.ContextMiddleware(routes.Gallery)).Methods("GET")
	server.Router.HandleFunc("/{username}/Gallery/feed", server.ContextMiddleware(routes.GalleryFeed)).Methods("GET")

	// Account pages
	server.Router.HandleFunc("/account", server.ContextMiddleware(routes.Account)).Methods("GET")
	server.Router.HandleFunc("/account/", server.ContextMiddleware(routes.Account)).Methods("GET")
	server.Router.HandleFunc("/account/search/", server.ContextMiddleware(routes.Account)).Methods("GET")
	server.Router.HandleFunc("/account/go_pro", server.ContextMiddleware(routes.AccountGoPro)).Methods("GET")
	server.Router.HandleFunc("/account/settings", server.ContextMiddleware(routes.AccountSettings)).Methods("GET")
	server.Router.HandleFunc("/account/subscription", server.ContextMiddleware(routes.AccountSubscription)).Methods("GET")
	server.Router.HandleFunc("/account/reset_api_key", server.ContextMiddleware(routes.AccountResetApiKey)).Methods("GET")

	// AJAX pages
	server.Router.HandleFunc("/ajax/move_dialog/", server.ContextMiddleware(routes.MoveDialog)).Methods("GET")
	server.Router.HandleFunc("/ajax/move_upload", server.ContextMiddleware(routes.MoveUpload)).Methods("POST")
	server.Router.HandleFunc("/ajax/delete_upload", server.ContextMiddleware(routes.DeleteUpload)).Methods("POST")
	server.Router.HandleFunc("/ajax/change_password", server.ContextMiddleware(routes.ChangePassword)).Methods("POST")
	server.Router.HandleFunc("/ajax/default_puush_pool", server.ContextMiddleware(routes.UpdateDefaultPool)).Methods("POST")
	server.Router.HandleFunc("/ajax/confirm_username", server.ContextMiddleware(routes.CheckUsername)).Methods("POST")
	server.Router.HandleFunc("/ajax/claim_username", server.ContextMiddleware(routes.ClaimUsername)).Methods("POST")
	server.Router.HandleFunc("/ajax/stopnagging", server.ContextMiddleware(routes.StopAskingAboutUsername)).Methods("POST")

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

func EnsureWebFolder() {
	// Ensure the web folder exists
	// If it doesn't, create it and populate it with default files
	for _, folder := range requiredFolders {
		if _, err := os.Stat(folder); !os.IsNotExist(err) {
			continue
		}
		fmt.Printf("Required folder '%s' does not exist. Downloading from GitHub...\n", folder)

		// Download the folder from github
		err := DownloadDirectory(folder)
		if err != nil {
			log.Fatalf("Failed to download required folder %s: %v", folder, err)
		}
	}
}

func init() {
	// Download web folder if it doesn't exist
	EnsureWebFolder()

	// Create .env file if it doesn't exist
	err := CreateDefaultEnvironment()
	if err != nil {
		log.Fatalf("Failed to create default .env file: %v", err)
	}
}

func main() {
	state := app.NewState()
	if state == nil {
		return
	}
	defer state.Close()

	server := app.NewServer(
		state.Config.Web.Host,
		state.Config.Web.Port,
		"puush",
		state,
	)
	InitializeWebRoutes(server)
	InitializeApiRoutes(server)
	InitializeCdnRoutes(server)
	server.Serve()
}
