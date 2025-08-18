package release

import (
	"encoding/json"
	"net/http"
	"peanut/internal/cache"
	"peanut/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func GetLatestRelease(w http.ResponseWriter, r *http.Request) {
	repo := chi.URLParam(r, "repository")
	latest, err := cache.FindLatest(repo)
	if err != nil {
		log.Error().Err(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	release := repository.Release{
		Name:         latest.Name,
		Body:         latest.Body,
		TagName:      latest.TagName,
		IsDraft:      latest.IsDraft,
		IsPrerelease: latest.IsPrerelease,
		CreatedAt:    latest.CreatedAt,
		PublishedAt:  latest.PublishedAt,
		AuthorName:   latest.AuthorName,
	}

	err = json.NewEncoder(w).Encode(release)
	if err != nil {
		log.Error().Err(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}
