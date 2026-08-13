package server

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/Lekuruu/go-puush/internal/state"
)

const (
	httpReadHeaderTimeout = 5 * time.Second
	httpIdleTimeout       = 2 * time.Minute
	httpShutdownTimeout   = 30 * time.Second
)

// Server is the main struct that holds the state for an http server.
type Server struct {
	Host   string
	Port   int
	Name   string
	State  *state.State
	Router *http.ServeMux
	Logger *slog.Logger
}

func NewServer(host string, port int, name string, state *state.State) *Server {
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
		Router: http.NewServeMux(),
	}
}

// Handle registers an application handler using a standard-library route pattern.
func (server *Server) Handle(pattern string, handler func(*Context)) {
	server.Router.HandleFunc(pattern, server.ContextMiddleware(handler))
}

// Serve starts the HTTP server and listens for incoming requests.
func (server *Server) Serve(ctx context.Context) error {
	httpServer := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", server.Host, server.Port),
		Handler:           server.LoggingMiddleware(server.Router),
		ErrorLog:          slog.NewLogLogger(server.Logger.Handler(), slog.LevelError),
		ReadHeaderTimeout: httpReadHeaderTimeout,
		IdleTimeout:       httpIdleTimeout,
	}
	server.Logger.Info(
		"Listening for requests",
		"host", server.Host,
		"port", server.Port,
	)

	serveErrors := make(chan error, 1)
	go func() {
		serveErrors <- httpServer.ListenAndServe()
	}()

	select {
	case err := <-serveErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		server.Logger.Info("Shutting down server")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), httpShutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		closeErr := httpServer.Close()
		<-serveErrors
		return errors.Join(fmt.Errorf("server: failed to shut down gracefully: %w", err), closeErr)
	}

	err := <-serveErrors
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
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

func (rc *ResponseContext) Write(data []byte) (int, error) {
	rc.WriteImplicitStatus()
	return rc.w.Write(data)
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
	rc.FlushError()
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
		ctx := &Context{
			Response: w,
			Request:  r,
			State:    server.State,
			Logger:   server.Logger,
		}

		w.Header().Set("Server", server.Name)
		handler(ctx)
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
