package release

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"net/http"
	"peanut/internal/cache"
	"peanut/internal/github"
)

func GetLatestRelease(w http.ResponseWriter, r *http.Request) {
	repo := chi.URLParam(r, "repository")
	latest, err := cache.FindLatest(repo)
	if err != nil {
		log.Error().Err(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	release := github.Release{
		Name:        latest.Name,
		Body:        latest.Body,
		TagName:     latest.TagName,
		Draft:       latest.Draft,
		Prerelease:  latest.Prerelease,
		CreatedAt:   latest.CreatedAt,
		PublishedAt: latest.PublishedAt,
		AuthorName:  latest.Author.Name,
	}

	err = json.NewEncoder(w).Encode(release)
	if err != nil {
		log.Error().Err(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}
