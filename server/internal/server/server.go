package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"

	"go.ads.coffee/platform/server/internal/pipeline"
)

type Manager interface {
	Mount(router *chi.Mux)
}

type Server struct {
	config  Config
	srv     *http.Server
	manager Manager
}

func New(config Config, manager *pipeline.Manager) *Server {
	return &Server{
		config:  config,
		srv:     &http.Server{Addr: config.Port},
		manager: manager,
	}
}

func (s *Server) Start(ctx context.Context) error {
	router := chi.NewRouter()

	s.manager.Mount(router)

	http.Handle("/", router)

	ln, err := net.Listen("tcp", s.config.Port)
	if err != nil {
		return err
	}

	fmt.Println("Served at http://localhost" + s.config.Port)

	go func() {
		if err := s.srv.Serve(ln); err != nil {
			fmt.Println(err)
		}
	}()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
