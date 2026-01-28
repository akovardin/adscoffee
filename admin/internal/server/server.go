package server

import (
	"context"
	"net"
	"net/http"

	"github.com/qor5/admin/v3/presets"
	"github.com/qor5/x/v3/login"
	"go.uber.org/zap"
)

type Server struct {
	config Config
	logger *zap.Logger
	lb     *login.Builder
	pb     *presets.Builder
	srv    *http.Server
}

func New(
	config Config,
	logger *zap.Logger,
	lb *login.Builder,
	pb *presets.Builder,
) *Server {
	return &Server{
		config: config,
		logger: logger,
		lb:     lb,
		pb:     pb,
	}
}

func (s *Server) Serve() error {
	s.logger.Info("app server", zap.String("host", "http://localhost"+s.config.Port+"/admin"))

	mux := http.NewServeMux()
	mux.Handle("/", s.pb)
	// mux.Handle("/duAction1", http.HandlerFunc() func(w http.ResponseWriter, r *http.Request) {

	// })

	s.lb.Mount(mux)

	http.Handle("/", s.lb.Middleware()(mux))
	// http.Handle("/", mux)

	s.srv = &http.Server{Addr: s.config.Port}

	ln, err := net.Listen("tcp", s.srv.Addr)
	if err != nil {
		return err
	}

	go func() {
		if err := s.srv.Serve(ln); err != nil {
			s.logger.Error("app server", zap.Error(err))
		}
	}()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
