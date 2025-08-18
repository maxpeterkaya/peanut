package cache

import (
	"errors"
	"fmt"
	"peanut/internal/config"
	"peanut/internal/database"
	gh "peanut/internal/github"
	"peanut/internal/repository"
	"peanut/internal/services"
	"slices"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

var (
	Repositories []Repo
)

type Repo struct {
	Owner         string
	Token         string
	Name          string
	LatestRelease *repository.Release
	Releases      []repository.Release
	Assets        []repository.Asset
}

func Init() error {
	if config.Config.Common.EnableDatabase {
		repos, err := database.Queries.ListRepositories(database.CTX)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to get repositories to cache")
			return err
		}

		for _, repo := range repos {
			fullRepo := fmt.Sprintf("%v/%v", repo.Owner.String, repo.Name.String)

			releases, err := gh.GetReleases(repo.Owner.String, repo.Name.String)
			if err != nil {
				log.Error().Err(err).Str("repo", fullRepo).Msg("Failed to get releases")
				return err
			}

			for _, release := range releases {
				createRelease, err := services.CreateRelease(*release, fullRepo, repo)
				if err != nil {
					return err
				}

				for _, asset := range release.Assets {
					_, err := services.CreateAsset(*asset, fullRepo, *createRelease)
					if err != nil {
						return err
					}
				}
			}
		}

		return nil
	}

	for _, repo := range config.Config.Github.Repositories {
		releases, err := gh.GetReleases(config.Config.Github.Owner, repo)
		if err != nil {
			log.Error().Err(err).Str("repo", repo).Msg("Failed to get releases")
			return err
		}

		authorName := fmt.Sprintf("%v", releases[0].Author.Name)

		data := Repo{
			Owner: config.Config.Github.Owner,
			Token: config.Config.Github.Token,
			Name:  repo,
			LatestRelease: &repository.Release{
				ID:              int32(*releases[0].ID),
				GithubID:        int32(*releases[0].ID),
				TagName:         pgtype.Text{String: *releases[0].TagName},
				Name:            pgtype.Text{Valid: true, String: *releases[0].TagName},
				Body:            pgtype.Text{Valid: true, String: *releases[0].Body},
				IsDraft:         pgtype.Bool{Valid: true, Bool: *releases[0].Draft},
				IsPrerelease:    pgtype.Bool{Valid: true, Bool: *releases[0].Prerelease},
				PublishedAt:     pgtype.Timestamp{Time: releases[0].PublishedAt.Time, Valid: true},
				AuthorName:      pgtype.Text{Valid: true, String: authorName},
				AuthorID:        pgtype.Text{Valid: true, String: strconv.FormatInt(*releases[0].Author.ID, 10)},
				AuthorAvatarUrl: pgtype.Text{Valid: true, String: *releases[0].Author.AvatarURL},
				RepositoryID:    pgtype.Int4{Valid: true, Int32: int32(0)},
			},
			Releases: services.GithubReleasesToRepoReleases(releases),
			Assets:   services.GithubAssetsToRepoAssets(releases),
		}

		if len(Repositories) == 0 {
			Repositories = []Repo{data}
		}

		index := slices.IndexFunc(Repositories, func(r Repo) bool {
			return r.Name == repo
		})

		if index == -1 {
			Repositories = append(Repositories, data)
		} else {
			Repositories[index] = data
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

func FindVersion(repo string, version string) (*repository.Release, error) {
	if config.Config.Common.EnableDatabase {
		rp, err := services.SearchRepo(repo)
		if err != nil {
			log.Error().Err(err).Str("repo", repo).Str("version", version).Msg("Failed to search repo")
			return nil, err
		}

		release, err := database.Queries.GetReleaseVersion(database.CTX, repository.GetReleaseVersionParams{
			TagName:      pgtype.Text{Valid: true, String: version},
			RepositoryID: pgtype.Int4{Valid: true, Int32: rp.ID},
		})
		if err != nil {
			log.Error().Err(err).Str("repo", repo).Str("version", version).Msg("Failed to get release version")
			return nil, err
		}

		return &release, nil
	}

	for _, r := range Repositories {
		if r.Name == repo {
			for _, v := range r.Releases {
				if v.TagName.String == version {
					return &v, nil
				}
			}
		}
	}

	return nil, nil
}

func FindLatest(repo string) (*repository.Release, error) {
	if config.Config.Common.EnableDatabase {
		rp, err := services.SearchRepo(repo)
		if err != nil {
			log.Error().Err(err).Str("repo", repo).Str("version", "latest").Msg("Failed to search repo")
			return nil, err
		}

		release, err := database.Queries.GetLatestRelease(database.CTX, pgtype.Int4{Valid: true, Int32: rp.ID})
		if err != nil {
			log.Error().Err(err).Str("repo", repo).Str("version", "latest").Msg("Failed to get latest release")
		}

		return &release, nil
	}

	for _, r := range Repositories {
		if r.Name == repo {
			return r.LatestRelease, nil
		}
	}

	return nil, errors.New("no releases found")
}

func FindReleases(repo string) ([]repository.Release, error) {
	if config.Config.Common.EnableDatabase {
		rp, err := services.SearchRepo(repo)
		if err != nil {
			log.Error().Err(err).Str("repo", repo).Msg("Failed to search repo")
			return nil, err
		}

		releases, err := services.ListReleases(rp.ID)
		if err != nil {
			log.Error().Err(err).Str("repo", repo).Msg("Failed to list releases")
			return nil, err
		}

		return releases, nil
	}

	for _, r := range Repositories {
		if r.Name == repo {
			return r.Releases, nil
		}
	}

	return nil, errors.New("no releases found")
}

func FindAssets(releaseID int32) ([]repository.Asset, error) {
	if config.Config.Common.EnableDatabase {
		assets, err := services.GetAssets(releaseID)
		if err != nil {
			log.Trace().Err(err).Int("releaseID", int(releaseID)).Msg("Failed to get assets")
			return nil, err
		}

		return assets, nil
	}

	array := make([]repository.Asset, 0)

	for _, r := range Repositories {
		for _, asset := range r.Assets {
			if asset.ReleaseID.Int32 == releaseID {
				array = append(array, asset)
			}
		}
	}

	if len(array) == 0 {
		return nil, errors.New("no assets found")
	}

	return array, nil
}
