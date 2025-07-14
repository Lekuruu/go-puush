package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Lekuruu/go-puush/internal/app"
	"github.com/gorilla/mux"
)

// Server is the main struct that holds the state for the API server.
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

// Context is a struct that holds the request context for each API call.
type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	State    *app.State
}

// Serve starts the HTTP server and listens for incoming requests.
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

	err := http.ListenAndServe(bind, server.LoggingMiddleware(server.Router))
	if err != nil {
		slog.Error("Failed to start server", "error", err)
		return
	}
}

func (server *Server) InitializeRoutes() {
	server.Router.HandleFunc("/api/register", server.ContextMiddleware(PuushRegistration)).Methods("POST")
	server.Router.HandleFunc("/api/auth", server.ContextMiddleware(PuushAuthentication)).Methods("POST")
	server.Router.HandleFunc("/api/up", server.ContextMiddleware(PuushUpload)).Methods("POST")
	server.Router.HandleFunc("/api/del", server.ContextMiddleware(PuushDelete)).Methods("POST")
	server.Router.HandleFunc("/api/hist", server.ContextMiddleware(PuushHistory)).Methods("POST")
	server.Router.HandleFunc("/api/thumb", server.ContextMiddleware(PuushThumbnail)).Methods("POST")
	server.Router.HandleFunc("/api/oshi", server.ContextMiddleware(PuushErrorSubmission)).Methods("POST")
	server.Router.HandleFunc("/dl/puush-rss.xml", server.ContextMiddleware(PuushMacOSRssFeed)).Methods("GET")
	server.Router.HandleFunc("/dl/puush-win.txt", server.ContextMiddleware(PuushWindowsUpdate)).Methods("GET")
}

// ResponseContext is a wrapper around http.ResponseWriter that
// allows us to capture the status code of a response.
type ResponseContext struct {
	w http.ResponseWriter
	s int
}

func (rc *ResponseContext) Header() http.Header {
	return rc.w.Header()
}

func (rc *ResponseContext) Write(b []byte) (int, error) {
	return rc.w.Write(b)
}

func (rc *ResponseContext) WriteHeader(status int) {
	rc.s = status
	rc.w.WriteHeader(status)
}

func (rc *ResponseContext) Status() int {
	if rc.s == 0 {
		return http.StatusOK
	}
	return rc.s
}

// ContextMiddleware creates a new Context struct for each request.
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

// LoggingMiddleware logs the details of each request.
func (server *Server) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rc := &ResponseContext{w: w}
		start := time.Now()
		next.ServeHTTP(rc, r)
		time := time.Since(start)

		server.Logger.Info(
			"API Request",
			"method", r.Method,
			"uri", r.RequestURI,
			"remote", r.RemoteAddr,
			"user-agent", r.UserAgent(),
			"status", rc.Status(),
			"duration", time.String(),
		)
	})
}
