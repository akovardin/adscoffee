package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"

	"go.ads.coffee/platform/server/internal/pipeline"
)

const addr = ":9090"

type Manager interface {
	Mount(router *chi.Mux)
}

type Server struct {
	srv     *http.Server
	manager Manager
}

func New(manager *pipeline.Manager) *Server {
	return &Server{
		srv:     &http.Server{Addr: addr},
		manager: manager,
	}
}

func (s *Server) Start(ctx context.Context) error {
	router := chi.NewRouter()

	s.manager.Mount(router)

	http.Handle("/", router)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	fmt.Println("Served at http://localhost" + addr)

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
