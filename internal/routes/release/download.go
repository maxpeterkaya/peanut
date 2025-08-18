package release

import (
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"

	"net/http"
	"peanut/internal/helper"
	"peanut/internal/services"
	"strings"
)

func DownloadPlatform(w http.ResponseWriter, r *http.Request) {
	platform := chi.URLParam(r, "platform")
	repo := chi.URLParam(r, "repository")

	ua := helper.ParseUserAgent(platform)

	rp, err := services.SearchRepo(repo)
	if err != nil {
		log.Error().Err(err).Str("repo", repo).Str("platform", platform).Msg("Failed to search repo")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	release := helper.GetPlatformRelease(repo, ua.FirstPrediction)
	if release == nil {
		release = helper.GetPlatformRelease(repo, ua.SecondPrediction)
	}

	if release == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	helper.ProxyDownload(*rp, *release, w, r)
	return
}

func Download(w http.ResponseWriter, r *http.Request) {
	repo := chi.URLParam(r, "repository")

	ua := helper.ParseUserAgent(strings.ToLower(r.UserAgent()))

	rp, err := services.SearchRepo(repo)
	if err != nil {
		log.Error().Err(err).Str("repo", repo).Str("first_prediction", ua.FirstPrediction).Str("second_prediction", ua.SecondPrediction).Msg("Failed to search repo")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	release := helper.GetPlatformRelease(repo, ua.FirstPrediction)
	if release == nil {
		release = helper.GetPlatformRelease(repo, ua.SecondPrediction)
	}

	if release == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	helper.ProxyDownload(*rp, *release, w, r)
	return
}
