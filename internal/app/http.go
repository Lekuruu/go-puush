package app

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Server is the main struct that holds the state for an http server.
type Server struct {
	Host   string
	Port   int
	Name   string
	State  *State
	Router *mux.Router
	Logger *slog.Logger
}

func NewServer(host string, port int, name string, state *State) *Server {
	logger := slog.Default()
	if state != nil && state.Logger != nil {
		logger = state.Logger
	}

	return &Server{
		Host:   host,
		Port:   port,
		Name:   name,
		State:  state,
		Logger: logger.With("component", name),
		Router: mux.NewRouter(),
	}
}

// Context is a struct that holds the request context for each endpoint call.
type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	State    *State
	Logger   *slog.Logger
	Vars     map[string]string
}

func (ctx *Context) IP() string {
	return GetRequestIP(ctx.Request)
}

// Serve starts the HTTP server and listens for incoming requests.
func (server *Server) Serve() {
	bind := fmt.Sprintf(
		"%s:%d",
		server.Host,
		server.Port,
	)
	server.Logger.Info(
		"Listening for requests",
		"host", server.Host,
		"port", server.Port,
	)

	httpServer := &http.Server{
		Addr:     bind,
		Handler:  server.LoggingMiddleware(server.Router),
		ErrorLog: slog.NewLogLogger(server.Logger.Handler(), slog.LevelError),
	}
	err := httpServer.ListenAndServe()
	if err != nil {
		server.Logger.Error("HTTP server stopped", "error", err)
		return
	}
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
	rc.WriteImplicitStatus()
	return rc.w.Write(b)
}

func (rc *ResponseContext) WriteHeader(status int) {
	// Informational responses do not contain the final response status.
	if status >= 100 && status < 200 && status != http.StatusSwitchingProtocols {
		rc.w.WriteHeader(status)
		return
	}
	if rc.s != 0 {
		return
	}
	rc.s = status
	rc.w.WriteHeader(status)
}

func (rc *ResponseContext) Unwrap() http.ResponseWriter {
	return rc.w
}

func (rc *ResponseContext) Flush() {
	_ = rc.FlushError()
}

func (rc *ResponseContext) FlushError() error {
	rc.WriteImplicitStatus()
	return http.NewResponseController(rc.w).Flush()
}

func (rc *ResponseContext) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return http.NewResponseController(rc.w).Hijack()
}

func (rc *ResponseContext) Push(target string, options *http.PushOptions) error {
	writer := rc.w
	for {
		if pusher, ok := writer.(http.Pusher); ok {
			return pusher.Push(target, options)
		}
		unwrapper, ok := writer.(interface{ Unwrap() http.ResponseWriter })
		if !ok {
			return http.ErrNotSupported
		}
		writer = unwrapper.Unwrap()
	}
}

func (rc *ResponseContext) ReadFrom(reader io.Reader) (int64, error) {
	rc.WriteImplicitStatus()
	if readerFrom, ok := rc.w.(io.ReaderFrom); ok {
		return readerFrom.ReadFrom(reader)
	}

	// Hide ReadFrom from io.Copy to avoid recursively calling this method.
	return io.Copy(struct{ io.Writer }{rc}, reader)
}

func (rc *ResponseContext) WriteImplicitStatus() {
	if rc.s == 0 {
		rc.s = http.StatusOK
	}
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
			Logger:   server.Logger,
			Vars:     mux.Vars(r),
		}

		w.Header().Set("Server", server.Name)
		handler(context)
	}
}

// LoggingMiddleware logs the details of each request.
func (server *Server) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rc := &ResponseContext{w: w}
		start := time.Now()
		next.ServeHTTP(rc, r)
		server.Logger.Info(
			fmt.Sprintf("%s %s", r.Method, r.URL.EscapedPath()),
			"ip", GetRequestIP(r),
			"status", rc.Status(),
			"duration", time.Since(start).String(),
		)
	})
}
