package web

import (
	"context"
	"net"
	"net/http"
	"time"
)

type Port string

func ProvideServer(port Port) *Server {
	return &Server{
		port: port,
	}
}

type Server struct {
	port Port
}

func (s *Server) Start(ctx context.Context) error {
	hs := &http.Server{
		Addr:              net.JoinHostPort("", string(s.port)),
		Handler:           s.handler(),
		ReadHeaderTimeout: time.Second * 3,
	}
	return start(ctx, hs, time.Second*5)
}

func (s *Server) handler() http.Handler {
	mux := http.NewServeMux()
	return mux
}
