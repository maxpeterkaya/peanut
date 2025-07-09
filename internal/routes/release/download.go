package release

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"peanut/internal/config"
	"peanut/internal/helper"
	"strings"
)

func DownloadPlatform(w http.ResponseWriter, r *http.Request) {
	platform := chi.URLParam(r, "platform")
	repo := chi.URLParam(r, "repository")

	ua := helper.ParseUserAgent(platform)

	release := helper.GetPlatformRelease(repo, ua.FirstPrediction)
	if release == nil {
		release = helper.GetPlatformRelease(repo, ua.SecondPrediction)
	}

	if release == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	helper.ProxyDownload(*release, config.Config.Github.Token, w, r)
	return
}

func Download(w http.ResponseWriter, r *http.Request) {
	repo := chi.URLParam(r, "repository")

	ua := helper.ParseUserAgent(strings.ToLower(r.UserAgent()))

	release := helper.GetPlatformRelease(repo, ua.FirstPrediction)
	if release == nil {
		release = helper.GetPlatformRelease(repo, ua.SecondPrediction)
	}

	if release == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	helper.ProxyDownload(*release, config.Config.Github.Token, w, r)
	return
}
