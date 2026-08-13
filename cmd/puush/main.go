package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Lekuruu/go-puush/api"
	"github.com/Lekuruu/go-puush/cdn"
	"github.com/Lekuruu/go-puush/internal/server"
	"github.com/Lekuruu/go-puush/internal/state"
	"github.com/Lekuruu/go-puush/web/routes"
)

func InitializeCdnRoutes(instance *server.Server) {
	// Gallery pages
	galleryRouter := http.NewServeMux()
	galleryRouter.HandleFunc("GET /{username}/Gallery", instance.ContextMiddleware(routes.Gallery))
	galleryRouter.HandleFunc("GET /{username}/Gallery/feed", instance.ContextMiddleware(routes.GalleryFeed))

	cdnRouter := http.NewServeMux()
	cdnRouter.HandleFunc("GET /t/{identifier}", instance.ContextMiddleware(cdn.Thumbnail))
	cdnRouter.HandleFunc("GET /t/{pool}/{identifier}", instance.ContextMiddleware(cdn.Thumbnail))
	cdnRouter.HandleFunc("GET /{identifier}", instance.ContextMiddleware(cdn.Upload))
	cdnRouter.HandleFunc("GET /{pool}/{identifier}", instance.ContextMiddleware(cdn.Upload))

	instance.Router.Handle("GET /", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isGalleryRoute(r.URL.Path) {
			galleryRouter.ServeHTTP(w, r)
			return
		}
		cdnRouter.ServeHTTP(w, r)
	}))
}

func isGalleryRoute(path string) bool {
	parts := strings.Split(strings.Trim(path, "/"), "/")

	// Gallery routes always have "Gallery" as the second segment
	if len(parts) < 2 || parts[1] != "Gallery" {
		return false
	}

	// Match /{username}/Gallery and /{username}/Gallery/feed
	return len(parts) == 2 || len(parts) == 3 && parts[2] == "feed"
}

func InitializeApiRoutes(instance *server.Server) {
	instance.Handle("POST /api/register", api.PuushRegistration)
	instance.Handle("POST /api/auth", api.PuushAuthentication)
	instance.Handle("POST /api/up", api.PuushUpload)
	instance.Handle("POST /api/del", api.PuushDelete)
	instance.Handle("POST /api/hist", api.PuushHistory)
	instance.Handle("POST /api/thumb", api.PuushThumbnail)
	instance.Handle("POST /api/oshi", api.PuushErrorSubmission)
}

func InitializeWebRoutes(instance *server.Server) {
	// Initialize templates
	routes.InitializeTemplates()

	// Public pages
	instance.Handle("GET /{$}", routes.Home)
	instance.Handle("GET /faq", routes.Faq)
	instance.Handle("GET /about", routes.About)
	instance.Handle("GET /register", routes.Register)
	instance.Handle("GET /register/{$}", routes.Register)
	instance.Handle("POST /register/go", routes.PerformRegistration)
	instance.Handle("GET /register/verify", routes.PerformActivation)
	instance.Handle("GET /reset_password", routes.ResetPassword)
	instance.Handle("POST /reset_password", routes.RequestPasswordReset)
	instance.Handle("GET /reset_password/go", routes.ShowPasswordResetForm)
	instance.Handle("POST /reset_password/go", routes.PerformPasswordReset)
	instance.Handle("GET /tos", routes.TermsOfService)
	instance.Handle("GET /dmca", routes.Dmca)

	// Login pages
	instance.Handle("GET /login", routes.Login)
	instance.Handle("GET /login/{$}", routes.Login)
	instance.Handle("GET /login/go", routes.PerformLogin)
	instance.Handle("POST /login/go", routes.PerformLogin)
	instance.Handle("GET /login/go/{$}", routes.PerformLogin)
	instance.Handle("POST /login/go/{$}", routes.PerformLogin)
	instance.Handle("GET /login/retry", routes.Login)
	instance.Handle("GET /login/retry/{$}", routes.Login)
	instance.Handle("GET /logout", routes.Logout)
	instance.Handle("GET /logout/{$}", routes.Logout)

	// Account pages
	instance.Handle("GET /account", routes.Account)
	instance.Handle("GET /account/{$}", routes.Account)
	instance.Handle("GET /account/search/{$}", routes.Account)
	instance.Handle("GET /account/go_pro", routes.AccountGoPro)
	instance.Handle("GET /account/settings", routes.AccountSettings)
	instance.Handle("GET /account/subscription", routes.AccountSubscription)
	instance.Handle("GET /account/reset_api_key", routes.AccountResetApiKey)

	// AJAX pages
	instance.Handle("GET /ajax/move_dialog/{$}", routes.MoveDialog)
	instance.Handle("POST /ajax/move_upload", routes.MoveUpload)
	instance.Handle("POST /ajax/delete_upload", routes.DeleteUpload)
	instance.Handle("POST /ajax/change_password", routes.ChangePassword)
	instance.Handle("POST /ajax/default_puush_pool", routes.UpdateDefaultPool)
	instance.Handle("POST /ajax/confirm_username", routes.CheckUsername)
	instance.Handle("POST /ajax/claim_username", routes.ClaimUsername)
	instance.Handle("POST /ajax/stopnagging", routes.StopAskingAboutUsername)

	// Thumbnail page
	instance.Handle("GET /thumb/view/{identifier}", routes.Thumbnail)

	// Initialize static routes
	instance.Router.Handle("GET /dl/", http.StripPrefix("/dl/", http.FileServer(http.Dir("web/static/dl/"))))
	instance.Router.Handle("GET /js/", http.StripPrefix("/js/", http.FileServer(http.Dir("web/static/js/"))))
	instance.Router.Handle("GET /css/", http.StripPrefix("/css/", http.FileServer(http.Dir("web/static/css/"))))
	instance.Router.Handle("GET /img/", http.StripPrefix("/img/", http.FileServer(http.Dir("web/static/img/"))))
	instance.Router.HandleFunc("GET /robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/static/robots.txt")
	})
	instance.Router.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
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
		slog.Info("Required folder is missing; downloading it", "folder", folder)

		// Download the folder from github
		err := DownloadDirectory(folder)
		if err != nil {
			slog.Error("Failed to download required folder", "folder", folder, "error", err)
			os.Exit(1)
		}
	}
}

func init() {
	// Download web folder if it doesn't exist
	EnsureWebFolder()

	// Create .env file if it doesn't exist
	err := CreateDefaultEnvironment()
	if err != nil {
		slog.Error("Failed to create default environment file", "error", err)
		os.Exit(1)
	}
}

func main() {
	appState, err := state.NewState(".env")
	if err != nil {
		slog.Error("Failed to initialize application", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := appState.Close(); err != nil {
			slog.Error("Failed to close application state", "error", err)
		}
	}()

	httpServer := server.NewServer(
		appState.Config.Web.Host,
		appState.Config.Web.Port,
		"puush",
		appState,
	)
	InitializeWebRoutes(httpServer)
	InitializeApiRoutes(httpServer)
	InitializeCdnRoutes(httpServer)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := httpServer.Serve(ctx); err != nil {
		slog.Error("HTTP server stopped unexpectedly", "error", err)
	}
}
