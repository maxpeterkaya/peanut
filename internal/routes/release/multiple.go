package release

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"net/http"
	"peanut/internal/cache"
	"peanut/internal/github"
	"strconv"
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

	releases := make([]github.Release, 0)

	for _, release := range data {
		releases = append(releases, github.Release{
			Name:        release.Name,
			Body:        release.Body,
			TagName:     release.TagName,
			Draft:       release.Draft,
			Prerelease:  release.Prerelease,
			CreatedAt:   release.CreatedAt,
			PublishedAt: release.PublishedAt,
			AuthorName:  release.Author.Name,
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
