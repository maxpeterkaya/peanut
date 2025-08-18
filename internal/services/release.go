package services

import (
	"errors"
	"fmt"
	"peanut/internal/database"
	"peanut/internal/repository"
	"strconv"

	"github.com/google/go-github/v67/github"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

func CreateRelease(release github.RepositoryRelease, repoName string, repo repository.Repository) (*repository.Release, error) {
	find, err := SearchRelease(release, repoName)
	if err != nil {
		return nil, err
	}
	if find != nil {
		return find, nil
	}

	authorName := fmt.Sprintf("%v", release.Author.Name)

	create, err := database.Queries.CreateRelease(database.CTX, repository.CreateReleaseParams{
		Name:            pgtype.Text{Valid: true, String: *release.Name},
		TagName:         pgtype.Text{Valid: true, String: *release.TagName},
		Body:            pgtype.Text{Valid: true, String: *release.Body},
		IsDraft:         pgtype.Bool{Valid: true, Bool: *release.Draft},
		IsPrerelease:    pgtype.Bool{Valid: true, Bool: *release.Prerelease},
		AuthorName:      pgtype.Text{Valid: true, String: authorName},
		PublishedAt:     pgtype.Timestamp{Time: release.PublishedAt.Time, InfinityModifier: 0, Valid: true},
		AuthorID:        pgtype.Text{Valid: true, String: strconv.FormatInt(*release.Author.ID, 10)},
		AuthorAvatarUrl: pgtype.Text{Valid: true, String: *release.Author.AvatarURL},
		RepositoryID:    pgtype.Int4{Valid: true, Int32: repo.ID},
		GithubID:        int32(*release.ID),
	})
	if err != nil {
		log.Error().Err(err).Str("repo", repoName).Msg("Failed to create release")
		return nil, err
	}

	return &create, nil
}

func SearchRelease(release github.RepositoryRelease, repoName string) (*repository.Release, error) {
	find, err := database.Queries.GetGithubRelease(database.CTX, int32(*release.ID))
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		log.Error().Err(err).Str("repo", repoName).Msg("Failed to search for release")
		return nil, err
	} else if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return &find, nil
}

func ListReleases(repositoryID int32) ([]repository.Release, error) {
	find, err := database.Queries.ListReleases(database.CTX, pgtype.Int4{Int32: repositoryID, Valid: true})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		log.Error().Err(err).Int("repositoryID", int(repositoryID)).Msg("Failed to list releases")
		return nil, err
	}
	return find, nil
}

//func GetReleaseWithAssets(release github.RepositoryRelease, repoName string) (*repository.Release, []*repository.Asset, error) {
//
//}

func GithubReleasesToRepoReleases(releases []*github.RepositoryRelease) []repository.Release {
	array := make([]repository.Release, len(releases))

	for i, release := range releases {
		authorName := fmt.Sprintf("%v", release.Author.Name)

		array[i] = repository.Release{
			ID:              int32(*release.ID),
			GithubID:        int32(*release.ID),
			Name:            pgtype.Text{Valid: true, String: *release.TagName},
			Body:            pgtype.Text{Valid: true, String: *release.Body},
			IsDraft:         pgtype.Bool{Valid: true, Bool: *release.Draft},
			IsPrerelease:    pgtype.Bool{Valid: true, Bool: *release.Prerelease},
			PublishedAt:     pgtype.Timestamp{Time: release.PublishedAt.Time},
			AuthorName:      pgtype.Text{Valid: true, String: authorName},
			AuthorID:        pgtype.Text{Valid: true, String: strconv.FormatInt(*release.Author.ID, 10)},
			AuthorAvatarUrl: pgtype.Text{Valid: true, String: *release.Author.AvatarURL},
			RepositoryID:    pgtype.Int4{Valid: true, Int32: int32(0)},
		}
	}

	return array
}

func GithubAssetsToRepoAssets(releases []*github.RepositoryRelease) []repository.Asset {
	array := make([]repository.Asset, 0)

	for _, release := range releases {
		for _, asset := range release.Assets {
			array = append(array, repository.Asset{
				ID:            int32(*asset.ID),
				GithubID:      int32(*asset.ID),
				ApiUrl:        pgtype.Text{Valid: true, String: *asset.URL},
				Url:           pgtype.Text{Valid: true, String: *asset.BrowserDownloadURL},
				Name:          pgtype.Text{Valid: true, String: *asset.Name},
				ContentLength: pgtype.Int4{Valid: true, Int32: int32(*asset.Size)},
				DownloadCount: pgtype.Int4{Valid: true, Int32: int32(*asset.DownloadCount)},
				ViewCount:     pgtype.Int4{Valid: true, Int32: 0},
				CreatedAt:     pgtype.Timestamp{Valid: true, Time: release.PublishedAt.Time},
				UpdatedAt:     pgtype.Timestamp{Valid: true, Time: release.PublishedAt.Time},
				UploadedAt:    pgtype.Timestamp{Valid: true, Time: release.PublishedAt.Time},
				ReleaseID:     pgtype.Int4{Valid: true, Int32: int32(*release.ID)},
			})
		}
	}

	return array
}
