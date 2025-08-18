package release

import (
	"encoding/json"
	"net/http"
	"peanut/internal/cache"
	"peanut/internal/repository"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func GetVersionRelease(w http.ResponseWriter, r *http.Request) {
	tag := chi.URLParam(r, "tag")
	repo := chi.URLParam(r, "repository")

	data, err := cache.FindVersion(repo, tag)
	if err != nil {
		log.Error().Err(err).Str("tag", tag).Msg("GetVersionRelease")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	release := repository.Release{
		Name:         data.Name,
		Body:         data.Body,
		TagName:      data.TagName,
		IsDraft:      data.IsDraft,
		IsPrerelease: data.IsPrerelease,
		CreatedAt:    data.CreatedAt,
		PublishedAt:  data.PublishedAt,
		AuthorName:   data.AuthorName,
	}

	err = json.NewEncoder(w).Encode(release)
	if err != nil {
		log.Error().Err(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}
