package release

import (
	"encoding/json"
	"net/http"
	"peanut/internal/cache"
	"peanut/internal/repository"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func GetMultipleReleases(w http.ResponseWriter, r *http.Request) {
	amount, _ := strconv.Atoi(r.URL.Query().Get("amount"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	repo := chi.URLParam(r, "repository")

	rl, err := cache.FindReleases(repo)
	if err != nil {
		log.Error().Err(err).Msg("")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if len(rl) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	maxLength := len(rl)

	if amount == 0 || amount < 0 {
		amount = 20
	} else if amount > 100 {
		amount = 100
	}

	if amount > maxLength {
		amount = maxLength
	}

	data := rl[offset : offset+amount]

	releases := make([]repository.Release, 0)

	for _, release := range data {
		releases = append(releases, repository.Release{
			Name:         release.Name,
			Body:         release.Body,
			TagName:      release.TagName,
			IsDraft:      release.IsDraft,
			IsPrerelease: release.IsPrerelease,
			CreatedAt:    release.CreatedAt,
			PublishedAt:  release.PublishedAt,
			AuthorName:   release.AuthorName,
		})
	}

	err = json.NewEncoder(w).Encode(releases)
	if err != nil {
		log.Error().Err(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}
