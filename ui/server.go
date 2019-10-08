package ui

import (
	"context"
	"github.com/ns3777k/mailcage/pkg/httputils"
	"net/http"
	"time"

	"github.com/ns3777k/mailcage/pkg/httputils"

	"github.com/rs/zerolog"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
)

type Server struct {
	srv    *http.Server
	logger zerolog.Logger
	opts   *ServerOptions
}

type ServerOptions struct {
	ListenAddr string
	ForceAuth  bool
	Users      map[string]string
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func NewServer(opts *ServerOptions, logger zerolog.Logger) *Server {
	router := mux.NewRouter()
	srv := &http.Server{Addr: opts.ListenAddr, Handler: router}
	box := packr.New("assets", "./assets")

	router.HandleFunc("/healthcheck", healthcheck).Methods("GET")

	uiRouter := router.PathPrefix("/").Subrouter()
	uiRouter.Use(httputils.NewBasicAuthMiddleware(opts.Users, opts.ForceAuth))
	uiRouter.Path("/").Handler(http.FileServer(box))

	return &Server{srv: srv, logger: logger, opts: opts}
}

func (s *Server) Run(ctx context.Context) error {
	httpShutdownCh := make(chan struct{})

	go func() {
		<-ctx.Done()

		graceCtx, graceCancel := context.WithTimeout(context.Background(), time.Second*5)
		defer graceCancel()

		if err := s.srv.Shutdown(graceCtx); err != nil {
			s.logger.Error().Err(err).Msg("failed to shutdown gracefully")
		}

		httpShutdownCh <- struct{}{}
	}()

	s.logger.Info().Msg("starting to listen on " + s.opts.ListenAddr)

	err := s.srv.ListenAndServe()
	<-httpShutdownCh

	if err == http.ErrServerClosed {
		return nil
	}

	return err
}
