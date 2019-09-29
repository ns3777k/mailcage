package api

import (
    "context"
    "github.com/ns3777k/mailcage/api/v1"
    "github.com/ns3777k/mailcage/storage"
    "github.com/rs/zerolog"
    "net/http"
    "time"

    "github.com/gorilla/mux"
)

type Server struct {
    srv *http.Server
    logger zerolog.Logger
    opts *ServerOptions
}

type ServerOptions struct {
    ListenAddr string
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}

func NewServer(opts *ServerOptions, logger zerolog.Logger, storage storage.Storage) *Server {
    router := mux.NewRouter()
    srv := &http.Server{Addr: opts.ListenAddr, Handler: router}

    router.HandleFunc("/healthcheck", healthcheck).Methods("GET")
    v1.NewAPI(storage).RegisterRoutes(router.PathPrefix("/api/v1").Subrouter())

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

    s.logger.Info().Msg("starting to listen on "+s.opts.ListenAddr)

    err := s.srv.ListenAndServe()
    <-httpShutdownCh

    if err == http.ErrServerClosed {
        return nil
    }

    return err
}
