package release

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
	"peanut/internal/cache"
	"peanut/internal/github"
)

func GetLatestRelease(w http.ResponseWriter, r *http.Request) {
	release := github.Release{
		Name:        cache.LatestRelease.Name,
		Body:        cache.LatestRelease.Body,
		TagName:     cache.LatestRelease.TagName,
		Draft:       cache.LatestRelease.Draft,
		Prerelease:  cache.LatestRelease.Prerelease,
		CreatedAt:   cache.LatestRelease.CreatedAt,
		PublishedAt: cache.LatestRelease.PublishedAt,
		AuthorName:  cache.LatestRelease.Author.Name,
	}

	err := json.NewEncoder(w).Encode(release)
	if err != nil {
		log.Error().Err(err)
	}
	return
}
