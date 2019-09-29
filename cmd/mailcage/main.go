package main

import (
    "context"
    "fmt"
    "github.com/ns3777k/mailcage/api"
    "github.com/ns3777k/mailcage/smtp"
    "github.com/ns3777k/mailcage/storage"
    "github.com/ns3777k/mailcage/ui"
    "github.com/pkg/errors"
    "github.com/rs/zerolog"
    "golang.org/x/sync/errgroup"
    "gopkg.in/alecthomas/kingpin.v2"
    "os"
    "os/signal"
    "path/filepath"
    "syscall"
)

type Configuration struct {
    ListenAddr string
    DebugMode bool
    Hostname string
    SMTPListenAddr string
    UIListenAddr string
    Storage string
    StorageFSDirectory string
}

func handleSignals(cancel context.CancelFunc) {
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)

    <-c
    cancel()
}

func createLogger(debug bool) zerolog.Logger {
    l := zerolog.New(os.Stderr).With().Timestamp().Logger()

    if debug {
        l = l.Level(zerolog.DebugLevel)
    } else {
        l = l.Level(zerolog.InfoLevel)
    }

    return l
}

func createStorage(config *Configuration) storage.Storage {
    var ret storage.Storage

    switch config.Storage {
    case "memory":
        ret = storage.CreateMemoryStorage()
    case "fs":
        ret = storage.CreateFsStorage(config.StorageFSDirectory)
    }

    return ret
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    g, ctx := errgroup.WithContext(ctx)

    config := new(Configuration)
    app := kingpin.New(filepath.Base(os.Args[0]), "MailCage")
    app.HelpFlag.Short('h')

    app.Flag("api-bind-addr", "Address to listen on for api").
        Default("127.0.0.1:8080").
        Envar("API_BIND_ADDR").
        StringVar(&config.ListenAddr)

    app.Flag("smtp-bind-addr", "Address to listen on for smtp").
        Default("127.0.0.1:1025").
        Envar("SMTP_BIND_ADDR").
        StringVar(&config.SMTPListenAddr)

    app.Flag("ui-bind-addr", "Address to listen on for ui").
        Default("127.0.0.1:8025").
        Envar("UI_BIND_ADDR").
        StringVar(&config.UIListenAddr)

    app.Flag("hostname", "SMTP ehlo/helo hostname").
        Default("mailcage.example").
        Envar("HOSTNAME").
        StringVar(&config.Hostname)

    app.Flag("debug", "More verbose logging").
        Envar("DEBUG").
        BoolVar(&config.DebugMode)

    app.Flag("storage", "Type of storage to save messages (memory or fs)").
        Envar("STORAGE").
        Default("memory").
        EnumVar(&config.Storage, "memory", "fs")

    app.Flag("storage-fs-dir", "Directory to create a database in (fs only)").
        Envar("STORAGE_FS_DIR").
        Default("/tmp").
        StringVar(&config.StorageFSDirectory)

    if _, err := app.Parse(os.Args[1:]); err != nil {
        fmt.Fprintln(os.Stderr, errors.Wrap(err, "error parsing commandline arguments"))
        app.Usage(os.Args[1:])
        os.Exit(2)
    }

    go handleSignals(cancel)

    logger := createLogger(config.DebugMode)
    s := createStorage(config)

    if err := s.Init(); err != nil {
        fmt.Fprintln(os.Stderr, errors.Wrap(err, "error setting up storage"))
        os.Exit(2)
    }

    defer s.Shutdown()

    g.Go(func() error {
        apiLogger := logger.With().Str("component", "api").Logger()
        apiServer := api.NewServer(&api.ServerOptions{ListenAddr: config.ListenAddr}, apiLogger, s)
        return apiServer.Run(ctx)
    })

    g.Go(func() error {
        smtpLogger := logger.With().Str("component", "smtp").Logger()
        sopts := &smtp.ServerOptions{Hostname: config.Hostname, ListenAddr: config.SMTPListenAddr}
        smtpServer := smtp.NewServer(sopts, smtpLogger, s)
        return smtpServer.Run(ctx)
    })

    g.Go(func() error {
        uiLogger := logger.With().Str("component", "ui").Logger()
        uiServer := ui.NewServer(&ui.ServerOptions{ListenAddr: config.UIListenAddr}, uiLogger)
        return uiServer.Run(ctx)
    })

    if err := g.Wait(); err != nil {
        logger.Err(err).Msg("failed to wait on goroutines")
    }
}
