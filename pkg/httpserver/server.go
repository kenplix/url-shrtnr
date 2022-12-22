package httpserver

import (
	"context"
	"net"
	"net/http"
	"time"
)

// Server -.
type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

// New -.
func New(handler http.Handler, options ...Option) *Server {
	s := Server{
		server: &http.Server{
			Addr:              net.JoinHostPort(defaultHost, defaultPort),
			Handler:           handler,
			ReadTimeout:       defaultReadTimeout,
			ReadHeaderTimeout: defaultReadHeaderTimeout,
			WriteTimeout:      defaultWriteTimeout,
			IdleTimeout:       defaultIdleTimeout,
		},
		notify:          make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
	}
	Preset(options...).apply(&s)

	return &s
}

// Start -.
func (s *Server) Start() string {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()

	return s.server.Addr
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
