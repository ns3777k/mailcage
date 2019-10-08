package main

import (
    "bufio"
    "context"
    "fmt"
    "io"
    "os"
    "os/signal"
    "path/filepath"
    "strings"
    "syscall"

    "github.com/ns3777k/mailcage/api"
    "github.com/ns3777k/mailcage/smtp"
    "github.com/ns3777k/mailcage/storage"
    "github.com/ns3777k/mailcage/ui"
    "github.com/pkg/errors"
    "github.com/rs/zerolog"
    "golang.org/x/sync/errgroup"
    "gopkg.in/alecthomas/kingpin.v2"

    _ "github.com/mattn/go-sqlite3"
)

type Configuration struct {
	ListenAddr             string
	DebugMode              bool
	Hostname               string
	SMTPListenAddr         string
	UIListenAddr           string
	AuthFilePath           string
	Storage                string
	StorageSQLiteDirectory string
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

func createStorage(config *Configuration) (storage.Storage, error) {
	var ret storage.Storage
	var err error

	switch config.Storage {
	case "memory":
		ret = storage.CreateMemoryStorage()
	case "sqlite":
		ret = storage.CreateSQLiteStorage(config.StorageSQLiteDirectory)
    default:
        err = errors.New("storage not found")
	}

	return ret, err
}

func parseAuthFile(filename string) (map[string]string, error) {
    users := make(map[string]string)

    if filename == "" {
        return users, nil
    }

    file, err := os.Open(filename)
    if err != nil {
        return users, err
    }

    defer file.Close()

    reader := bufio.NewReader(file)
    counter := 0

    for {
        counter++
        l, err := reader.ReadString('\n')
        l = strings.TrimSpace(l)

        if len(l) > 0 {
            auth := strings.Split(l, ":")
            if len(auth) != 2 {
                return users, errors.Errorf("invalid auth format at line: %d", counter)
            }

            users[auth[0]] = auth[1]
        }

        if err == io.EOF {
            break
        }

        if err != nil {
            return users, err
        }
    }

    return users, nil
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
		EnumVar(&config.Storage, "memory", "sqlite")

	app.Flag("storage-sqlite-dir", "Directory to create a database in (fs only)").
		Envar("STORAGE_FS_DIR").
		Default("/tmp").
		StringVar(&config.StorageSQLiteDirectory)

	app.Flag("auth-file", "Path to auth file").
	    Default("").
	    Envar("AUTH_FILE").
	    StringVar(&config.AuthFilePath)

	if _, err := app.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "error parsing commandline arguments"))
		app.Usage(os.Args[1:])
		os.Exit(2)
	}

	go handleSignals(cancel)

	logger := createLogger(config.DebugMode)

    users, err := parseAuthFile(config.AuthFilePath)
    if err != nil {
        logger.Fatal().Err(err).Msg("error parsing auth file")
    }

	s, err := createStorage(config)
	if err != nil {
        logger.Fatal().Err(err).Msg("error creating storage")
    }

	if err := s.Init(); err != nil {
        logger.Fatal().Err(err).Msg("error setting up storage")
	}

	defer func() {
		if err := s.Shutdown(); err != nil {
			logger.Err(err).Msg("storage shutdown failed")
		}
	}()

	g.Go(func() error {
		apiLogger := logger.With().Str("component", "api").Logger()
		apiOptions := &api.ServerOptions{
		    ListenAddr: config.ListenAddr,
		    ForceAuth: len(config.AuthFilePath) > 0,
            Users: users,
		}

		apiServer := api.NewServer(apiOptions, apiLogger, s)
		return apiServer.Run(ctx)
	})

	g.Go(func() error {
		smtpLogger := logger.With().Str("component", "smtp").Logger()
		smtpOptions := &smtp.ServerOptions{Hostname: config.Hostname, ListenAddr: config.SMTPListenAddr}
		smtpServer := smtp.NewServer(smtpOptions, smtpLogger, s)
		return smtpServer.Run(ctx)
	})

	g.Go(func() error {
		uiLogger := logger.With().Str("component", "ui").Logger()
		uiOptions := &ui.ServerOptions{
		    ListenAddr: config.UIListenAddr,
            ForceAuth: len(config.AuthFilePath) > 0,
            Users: users,
		}

		uiServer := ui.NewServer(uiOptions, uiLogger)
		return uiServer.Run(ctx)
	})

	if err := g.Wait(); err != nil {
		logger.Err(err).Msg("failed to wait on goroutines")
	}
}
