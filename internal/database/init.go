package database

import (
	"context"
	"peanut/internal/config"
	"peanut/internal/repository"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
)

var (
	DB      *pgxpool.Pool
	Queries *repository.Queries
	CTX     context.Context
)

func Init() error {
	CTX = context.Background()

	conf, err := pgxpool.ParseConfig(config.Config.Database.URL())
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse database config")
	}

	conf.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	conn, err := pgxpool.NewWithConfig(CTX, conf)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to database")
		return err
	}

	DB = conn
	Queries = repository.New(conn)

	return nil
}

func Close() {
	log.Info().Msg("closing database connection...")
	DB.Close()
}
