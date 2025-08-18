package services

import (
	"errors"
	"peanut/internal/database"
	"peanut/internal/repository"

	"github.com/google/go-github/v67/github"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rs/zerolog/log"
)

func CreateAsset(asset github.ReleaseAsset, repoName string, release repository.Release) (*repository.Asset, error) {
	find, err := SearchAsset(asset, repoName)
	if err != nil {
		return nil, err
	}
	if find != nil {
		return find, nil
	}

	create, err := database.Queries.CreateAsset(database.CTX, repository.CreateAssetParams{
		ApiUrl:        pgtype.Text{Valid: true, String: *asset.URL},
		Url:           pgtype.Text{Valid: true, String: *asset.BrowserDownloadURL},
		Name:          pgtype.Text{Valid: true, String: *asset.Name},
		ContentLength: pgtype.Int4{Valid: true, Int32: int32(*asset.Size)},
		DownloadCount: pgtype.Int4{Valid: true, Int32: int32(*asset.DownloadCount)},
		UploadedAt:    pgtype.Timestamp{Time: release.PublishedAt.Time, InfinityModifier: 0, Valid: true},
		ViewCount:     pgtype.Int4{Valid: true, Int32: 0},
		ReleaseID:     pgtype.Int4{Valid: true, Int32: release.ID},
		GithubID:      int32(*asset.ID),
	})
	if err != nil {
		log.Error().Err(err).Str("repo", repoName).Msg("Failed to create asset")
		return nil, err
	}

	return &create, nil
}

func SearchAsset(asset github.ReleaseAsset, repoName string) (*repository.Asset, error) {
	find, err := database.Queries.GetGithubAsset(database.CTX, int32(*asset.ID))
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		log.Error().Err(err).Str("repo", repoName).Msg("Failed to search for asset")
		return nil, err
	} else if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return &find, nil
}

func GetAssets(releaseID int32) ([]repository.Asset, error) {
	assets, err := database.Queries.ListReleaseAssets(database.CTX, pgtype.Int4{Valid: true, Int32: releaseID})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		log.Error().Err(err).Int("releaseID", int(releaseID)).Msg("Failed to list assets")
		return nil, err
	} else if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	return assets, nil
}
