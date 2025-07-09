package cache

import (
	"errors"
	"github.com/google/go-github/v67/github"
	"github.com/rs/zerolog/log"
	"peanut/internal/config"
	gh "peanut/internal/github"
	"slices"
)

var (
	Repositories []Repo
)

type Repo struct {
	Owner         string
	Token         string
	Name          string
	LatestRelease *github.RepositoryRelease
	Releases      []*github.RepositoryRelease
}

func Init() error {
	for _, repo := range config.Config.Github.Repositories {
		releases, err := gh.GetReleases(config.Config.Github.Owner, repo)
		if err != nil {
			log.Error().Err(err).Str("repo", repo).Msg("Failed to get releases")
			return err
		}

		if len(Repositories) == 0 {
			Repositories = []Repo{
				{
					Owner:         config.Config.Github.Owner,
					Token:         config.Config.Github.Token,
					Name:          repo,
					LatestRelease: releases[0],
					Releases:      releases,
				},
			}
		}

		index := slices.IndexFunc(Repositories, func(r Repo) bool {
			return r.Name == repo
		})

		if index == -1 {
			Repositories = append(Repositories, Repo{
				Owner:         config.Config.Github.Owner,
				Token:         config.Config.Github.Token,
				Name:          repo,
				LatestRelease: releases[0],
				Releases:      releases,
			})
		} else {
			Repositories[index] = Repo{
				Owner:         config.Config.Github.Owner,
				Token:         config.Config.Github.Token,
				Name:          repo,
				LatestRelease: releases[0],
				Releases:      releases,
			}
		}
	}

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

func FindVersion(repo string, version string) (*github.RepositoryRelease, error) {
	for _, r := range Repositories {
		if r.Name == repo {
			for _, v := range r.Releases {
				if *v.Name == version {
					return v, nil
				}
			}
		}
	}

	return nil, nil
}

func FindLatest(repo string) (*github.RepositoryRelease, error) {
	for _, r := range Repositories {
		if r.Name == repo {
			return r.LatestRelease, nil
		}
	}

	return nil, errors.New("no releases found")
}

func FindReleases(repo string) ([]*github.RepositoryRelease, error) {
	for _, r := range Repositories {
		if r.Name == repo {
			return r.Releases, nil
		}
	}

	return nil, errors.New("no releases found")
}
