package httpserver

import (
	"net"
	"strconv"
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

func SetPort(port uint16) Option {
	return optionFunc(func(s *Server) {
		p := strconv.FormatUint(uint64(port), 10)
		s.server.Addr = net.JoinHostPort("", p)
	})
}

func SetReadTimeout(timeout time.Duration) Option {
	return optionFunc(func(s *Server) {
		s.server.ReadTimeout = timeout
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
