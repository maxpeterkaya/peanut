package cache

import (
	"github.com/google/go-github/v67/github"
	"github.com/rs/zerolog/log"
	"peanut/internal/config"
	gh "peanut/internal/github"
)

var (
	LatestRelease *github.RepositoryRelease
	Releases      []*github.RepositoryRelease
)

func Init() error {
	releases, _, err := gh.GHClient.Repositories.ListReleases(gh.GHCtx, config.Config.Github.RepoOwner, config.Config.Github.Repository, &github.ListOptions{PerPage: 100})
	if err != nil {
		log.Error().Err(err).Msg("Error getting releases")
		return err
	}

	Releases = releases
	LatestRelease = releases[0]

	return nil
}

func Refresh() {
	log.Info().Msg("Refreshing cache...")

	err := Init()
	if err != nil {
		log.Error().Err(err).Msg("Error refreshing cache")
		return
	}

	log.Info().Msg("Refreshed cache.")
}

func FindVersion(version string) (*github.RepositoryRelease, error) {
	for _, release := range Releases {
		if *release.TagName == version {
			return release, nil
		}
	}

	return nil, nil
}
