package main

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"peanut/internal/buildinfo"
	"peanut/internal/cache"
	"peanut/internal/config"
	"peanut/internal/cron"
	"peanut/internal/database"
	"peanut/internal/github"
	"peanut/internal/metrics"
	"peanut/internal/routes/release"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/maxpeterkaya/peanut/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	chiprometheus "github.com/toshi0607/chi-prometheus"
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

		log.Info().Msg("seeding database...")
		err = database.SeedDB()
		if err != nil {
			return
		}
		log.Info().Msg("seeding complete")
	}

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
	r.Use(Logger(&log.Logger))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/health"))

	// Define routes
	if config.Config.Common.EnablePrometheus {
		if config.Config.Common.EnableGithubMetrics {
			prometheus.MustRegister(
				metrics.TotalAssets,
				metrics.TotalRepositories,
				metrics.TotalDownloads,
				metrics.TotalReleases,
				metrics.TotalRequests,
			)
		}

		m := chiprometheus.New("peanut")
		m.MustRegisterDefault()

		r.Use(m.Handler)

		r.Handle("/metrics", basicAuth(promhttp.Handler(), "prometheus", config.Config.Authorization.PrometheusToken))
	}
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

		log.Info().Msg("shutting down server...")

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

func Logger(log *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			t1 := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			t2 := time.Now()

			if rec := recover(); rec != nil {
				log.Error().
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Str("host", r.Host).
					Str("protocol", r.Proto).
					Str("referrer", r.Referer()).
					Int("status", ww.Status()).
					Str("ip", r.RemoteAddr).
					Str("user_agent", r.UserAgent()).
					Int64("bytes_in", r.ContentLength).
					Str("bytes_in_human", common.ByteCountSI(r.ContentLength)).
					Int("bytes_out", ww.BytesWritten()).
					Str("bytes_out_human", common.ByteCountSI(int64(ww.BytesWritten()))).
					Dur("latency", t2.Sub(t1)).
					Str("latency_human", t2.Sub(t1).String()).
					Interface("recovered", rec).
					Msg("panic recovered")
				http.Error(ww, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			log.Info().
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Str("host", r.Host).
				Str("protocol", r.Proto).
				Str("referrer", r.Referer()).
				Int("status", ww.Status()).
				Str("ip", r.RemoteAddr).
				Str("user_agent", r.UserAgent()).
				Int64("bytes_in", r.ContentLength).
				Str("bytes_in_human", common.ByteCountSI(r.ContentLength)).
				Int("bytes_out", ww.BytesWritten()).
				Str("bytes_out_human", common.ByteCountSI(int64(ww.BytesWritten()))).
				Dur("latency", t2.Sub(t1)).
				Str("latency_human", t2.Sub(t1).String()).
				Msg("api")
		}
		return http.HandlerFunc(fn)
	}
}
