package main

import (
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"peanut/internal/cache"
	"peanut/internal/config"
	"peanut/internal/cron"
	"peanut/internal/github"
	"peanut/internal/routes/release"
	"syscall"
	"time"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func init() {
	log.Info().Str("version", version).Str("commit", commit).Str("date", date).Msg("")

	if err := godotenv.Load(); err != nil {
		log.Error().Err(err).Msg("Error loading .env file")
	}

	log.Info().Msg("initializing config...")
	err := config.Init()
	if err != nil {
		return
	}
	log.Info().Msg("initialized config.")

	log.Info().Msg("initializing GitHub client...")
	err = github.Init()
	if err != nil {
		return
	}
	log.Info().Msg("initialized github client.")

	log.Info().Msg("populating cache...")
	err = cache.Init()
	if err != nil {
		return
	}
	log.Info().Msg("populated cache.")

	log.Info().Msg("initializing cron...")
	err = cron.Init()
	if err != nil {
		return
	}
	log.Info().Msg("initialized cron.")
}

func main() {
	r := chi.NewRouter()

	// Initialize middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Define routes
	r.Route("/release", func(r chi.Router) {
		r.Get("/latest", release.GetLatestRelease)
		r.Get("/{tag}", release.GetVersionRelease)
		r.Get("/multiple", release.GetMultipleReleases)
	})
	r.Route("/download", func(r chi.Router) {
		r.Get("/{platform}", release.DownloadPlatform)
		r.Get("/", release.Download)
	})

	server := &http.Server{Addr: "0.0.0.0:3000", Handler: r}
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatal().Msg("graceful shutdown timed out.. forcing exit.")
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal().Err(err)
		}
		serverStopCtx()
	}()

	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err)
	}

	<-serverCtx.Done()
}
