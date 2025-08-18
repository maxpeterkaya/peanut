package database

import (
	"peanut/internal/config"
	"peanut/internal/repository"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/maxpeterkaya/peanut/common"
	"github.com/rs/zerolog/log"
)

func SeedDB() error {
	user, err := createUser()
	if err != nil {
		return err
	}
	if user == nil {
		return nil
	}

	err = createRepo(user)

	return nil
}

func createUser() (*repository.User, error) {
	users, err := Queries.ListUsers(CTX)
	if err != nil {
		log.Error().Err(err).Str("task", "seed").Msg("failed to list users")
		return nil, err
	}

	if len(users) == 0 {
		create, err := Queries.CreateUser(CTX, repository.CreateUserParams{
			Username:    pgtype.Text{Valid: true, String: config.Config.Admin.User},
			DisplayName: pgtype.Text{Valid: true, String: config.Config.Admin.User},
			PassHash:    pgtype.Text{Valid: true, String: common.SHAHash(config.Config.Admin.Pass)},
		})
		if err != nil {
			log.Error().Err(err).Str("task", "seed").Msg("failed to create user")
			return nil, err
		}

		log.Info().Str("task", "seed").Str("user", create.Username.String).Msg("admin user created")
		return &create, nil
	}

	if len(users) == 1 {
		return &users[0], nil
	}

	log.Info().Str("task", "seed").Msg("users already exists")
	return nil, nil
}

func createRepo(user *repository.User) error {
	repos, err := Queries.ListRepositories(CTX)
	if err != nil {
		log.Error().Err(err).Str("task", "seed").Msg("failed to list repositories")
		return err
	}

	if len(repos) == 0 {
		for _, repo := range config.Config.Github.Repositories {
			_, err := Queries.CreateRepository(CTX, repository.CreateRepositoryParams{
				UserID:    pgtype.Int4{Valid: true, Int32: user.ID},
				Owner:     pgtype.Text{Valid: true, String: config.Config.Github.Owner},
				Name:      pgtype.Text{Valid: true, String: repo},
				Token:     pgtype.Text{Valid: true, String: config.Config.Github.Token},
				IsPrivate: pgtype.Bool{Valid: true, Bool: true},
				GithubID:  0,
			})
			if err != nil {
				log.Error().Err(err).Str("task", "seed").Msg("failed to create repository")
				return err
			}
		}
	}

	return nil
}
