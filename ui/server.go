package ui

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gobuffalo/packr/v2"

	"github.com/gorilla/mux"
	"github.com/ns3777k/mailcage/pkg/httputils"
	"github.com/rs/zerolog"
)

type Server struct {
	srv    *http.Server
	logger zerolog.Logger
	opts   *ServerOptions
}

type ServerOptions struct {
	ListenAddr        string
	ForceAuth         bool
	Users             map[string]string
	UIAssetsProxyAddr string
	APIProxyAddr      string
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func NewServer(opts *ServerOptions, logger zerolog.Logger) *Server {
	router := mux.NewRouter()
	srv := &http.Server{Addr: opts.ListenAddr, Handler: router}

	router.HandleFunc("/healthcheck", healthcheck).Methods("GET")

	uiRouter := router.PathPrefix("/").Subrouter()
	uiRouter.Use(httputils.NewBasicAuthMiddleware(opts.Users, opts.ForceAuth))

	apiURL, _ := url.Parse("http://" + opts.APIProxyAddr)
	apiProxy := httputil.NewSingleHostReverseProxy(apiURL)
	uiRouter.PathPrefix("/api").Handler(apiProxy)

	if opts.UIAssetsProxyAddr != "" {
		reactDevURL, err := url.Parse(opts.UIAssetsProxyAddr)
		if err != nil {
			panic(err)
		}
		uiProxy := httputil.NewSingleHostReverseProxy(reactDevURL)
		uiRouter.PathPrefix("/").Handler(uiProxy)
	} else {
		box := packr.New("assets", "./frontend/build")
		uiRouter.PathPrefix("/").Handler(http.FileServer(box))
	}

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
