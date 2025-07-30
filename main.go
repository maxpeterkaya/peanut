package main

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	chiprometheus "github.com/toshi0607/chi-prometheus"
	"net/http"
	"os"
	"os/signal"
	"peanut/internal/buildinfo"
	"peanut/internal/cache"
	"peanut/internal/config"
	"peanut/internal/cron"
	"peanut/internal/database"
	"peanut/internal/github"
	"peanut/internal/routes/release"
	"strings"
	"syscall"
	"time"
)

func init() {
	log.Info().Str("version", buildinfo.Version).Str("commit", buildinfo.Commit).Str("date", buildinfo.BuildTime).Msg("")

	if err := godotenv.Load(); err != nil {
		log.Error().Err(err).Msg("Error loading .env file")
	}

	log.Info().Msg("initializing config...")
	err := config.Init()
	if err != nil {
		os.Exit(1)
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

	if config.Config.Common.EnableDatabase {
		log.Info().Msg("initializing database...")
		err = database.Init()
		if err != nil {
			return
		}
		log.Info().Msg("initialized database.")

		log.Info().Msg("migrating database...")
		err = database.Migrate()
		if err != nil {
			return
		}
		log.Info().Msg("migrate complete")
	}
}

func main() {
	r := chi.NewRouter()

	m := chiprometheus.New("peanut")
	m.MustRegisterDefault()

	// Initialize middleware
	r.Use(m.Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/health"))

	// Define routes
	r.Handle("/metrics", basicAuth(promhttp.Handler(), "prometheus", config.Config.Authorization.PrometheusToken))
	// All routes that are related to repositories
	r.Route("/{repository}", func(r chi.Router) {
		r.Route("/release", func(r chi.Router) {
			r.Get("/latest", release.GetLatestRelease)
			r.Get("/{tag}", release.GetVersionRelease)
			r.Get("/multiple", release.GetMultipleReleases)
		})
		r.Route("/download", func(r chi.Router) {
			r.Get("/{platform}", release.DownloadPlatform)
			r.Get("/", release.Download)
		})
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
		database.Close()
		serverStopCtx()
	}()

	err := server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal().Err(err)
	}

	<-serverCtx.Done()
}

func basicAuth(h http.Handler, route, token string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ensure that authorization is actually enabled before proceeding
		if route == "prometheus" && !config.Config.Common.AuthPrometheus {
			h.ServeHTTP(w, r)
		} else if route == "github" && !config.Config.Common.AuthGithubMetrics {
			h.ServeHTTP(w, r)
		} else if route == "health" && !config.Config.Common.AuthHealthChek {
			h.ServeHTTP(w, r)
		}

		// Check authentication now...
		auth := r.Header.Get("Authorization")
		if !checkAuth(auth, token) {
			w.WriteHeader(401)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func checkAuth(auth string, token string) bool {
	prefix := "Basic "
	if !strings.HasPrefix(auth, prefix) {
		return false
	}
	payload, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return false
	}
	return string(payload) == token
}
