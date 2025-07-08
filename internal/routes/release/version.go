package release

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"net/http"
	"peanut/internal/cache"
	"peanut/internal/github"
)

func GetVersionRelease(w http.ResponseWriter, r *http.Request) {
	tag := chi.URLParam(r, "tag")

	data, err := cache.FindVersion(tag)
	if err != nil {
		log.Error().Err(err).Str("tag", tag).Msg("GetVersionRelease")
		return
	}

	release := github.Release{
		Name:        data.Name,
		Body:        data.Body,
		TagName:     data.TagName,
		Draft:       data.Draft,
		Prerelease:  data.Prerelease,
		CreatedAt:   data.CreatedAt,
		PublishedAt: data.PublishedAt,
		AuthorName:  data.Author.Name,
	}

	err = json.NewEncoder(w).Encode(release)
	if err != nil {
		log.Error().Err(err)
	}
	return
}
