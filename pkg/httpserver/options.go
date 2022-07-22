package httpserver

import (
	"net"
	"time"
)

// Option configures a Server.
type Option interface {
	apply(s *Server)
}

type optionFunc func(s *Server)

func (fn optionFunc) apply(s *Server) {
	fn(s)
}

// Preset turns a list of Option instances into an Option
func Preset(options ...Option) Option {
	return optionFunc(func(s *Server) {
		for _, option := range options {
			option.apply(s)
		}
	})
}

func SetPort(port string) Option {
	return optionFunc(func(s *Server) {
		s.server.Addr = net.JoinHostPort("", port)
	})
}

func SetReadTimeout(timeout time.Duration) Option {
	return optionFunc(func(s *Server) {
		s.server.ReadTimeout = timeout
	})
}

func SetReadHeaderTimeout(timeout time.Duration) Option {
	return optionFunc(func(s *Server) {
		s.server.ReadHeaderTimeout = timeout
	})
}

func SetIdleTimeout(timeout time.Duration) Option {
	return optionFunc(func(s *Server) {
		s.server.IdleTimeout = timeout
	})
}

func SetWriteTimeout(timeout time.Duration) Option {
	return optionFunc(func(s *Server) {
		s.server.WriteTimeout = timeout
	})
}

func SetShutdownTimeout(timeout time.Duration) Option {
	return optionFunc(func(s *Server) {
		s.shutdownTimeout = timeout
	})
}
