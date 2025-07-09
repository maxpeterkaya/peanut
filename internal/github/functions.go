package github

import (
	"context"
	"github.com/google/go-github/v67/github"
	"github.com/rs/zerolog/log"
	"peanut/internal/config"
)

var (
	GHClient *github.Client
	GHCtx    = context.Background()
)

func Init() error {
	if len(config.Config.Github.Token) > 0 {
		GHClient = github.NewClient(nil).WithAuthToken(config.Config.Github.Token)
	} else {
		GHClient = github.NewClient(nil)
	}

	return nil
}

func GetReleases(owner string, repo string) ([]*github.RepositoryRelease, error) {
	releases, _, err := GHClient.Repositories.ListReleases(GHCtx, owner, repo, &github.ListOptions{PerPage: 100})
	if err != nil {
		log.Error().Err(err).Msg("Error getting releases")
		return nil, err
	}

	return releases, nil
}
