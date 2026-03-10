package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"go.ads.coffee/platform/server/internal/pipeline"
	"go.ads.coffee/platform/server/static"
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

	router.Use(
		cors.Handler(cors.Options{
			// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
			// AllowedOrigins: []string{"https://*", "http://*"},
			AllowedOrigins: []string{"*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}),
	)

	s.manager.Mount(router)

	http.Handle("/", router)

	fs := http.FileServer(http.FS(static.FS))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

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
