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

	ua := helper.ParseUserAgent(platform)

	release := helper.GetPlatformRelease(ua.FirstPrediction)
	if release == nil {
		release = helper.GetPlatformRelease(ua.SecondPrediction)
	}

	if release == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	helper.ProxyDownload(*release, config.Config.GHToken, w, r)
	return
}

func Download(w http.ResponseWriter, r *http.Request) {
	ua := helper.ParseUserAgent(strings.ToLower(r.UserAgent()))

	release := helper.GetPlatformRelease(ua.FirstPrediction)
	if release == nil {
		release = helper.GetPlatformRelease(ua.SecondPrediction)
	}

	if release == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	helper.ProxyDownload(*release, config.Config.GHToken, w, r)
	return
}
