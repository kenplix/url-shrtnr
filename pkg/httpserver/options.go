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

func SetAddr(host, port string) Option {
	return optionFunc(func(s *Server) {
		if port == "" {
			port = defaultPort
		}

		s.server.Addr = net.JoinHostPort(host, port)
	})
}

func SetReadTimeout(timeout time.Duration) Option {
	return optionFunc(func(s *Server) {
		if timeout == 0 {
			timeout = defaultReadTimeout
		}

		s.server.ReadTimeout = timeout
	})
}

func SetReadHeaderTimeout(timeout time.Duration) Option {
	return optionFunc(func(s *Server) {
		if timeout == 0 {
			timeout = defaultReadHeaderTimeout
		}

		s.server.ReadHeaderTimeout = timeout
	})
}

func SetWriteTimeout(timeout time.Duration) Option {
	return optionFunc(func(s *Server) {
		if timeout == 0 {
			timeout = defaultWriteTimeout
		}

		s.server.WriteTimeout = timeout
	})
}

func SetIdleTimeout(timeout time.Duration) Option {
	return optionFunc(func(s *Server) {
		if timeout == 0 {
			timeout = defaultIdleTimeout
		}

		s.server.IdleTimeout = timeout
	})
}

func SetShutdownTimeout(timeout time.Duration) Option {
	return optionFunc(func(s *Server) {
		if timeout == 0 {
			timeout = defaultShutdownTimeout
		}

		s.shutdownTimeout = timeout
	})
}
