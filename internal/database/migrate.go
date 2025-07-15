package database

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog/log"
	"path/filepath"
	"peanut/internal/config"
)

func Migrate() error {
	path, _ := filepath.Abs("./db/migrations")

	m, err := migrate.New(
		fmt.Sprintf("file://%s", path),
		config.Config.Database.URL())
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to database")
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Error().Err(err).Msg("failed to run migrations")
		return err
	}

	err, err2 := m.Close()
	if err != nil {
		log.Error().Err(err).Msg("failed to close source")
		return err
	}
	if err2 != nil {
		log.Error().Err(err2).Msg("failed to close database")
		return err2
	}

	return nil
}
