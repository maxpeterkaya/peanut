package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
	"peanut/internal/config"
	"peanut/internal/repository"
)

var (
	DB      *pgx.Conn
	Queries *repository.Queries
	CTX     context.Context
)

func Init() error {
	CTX = context.Background()

	conn, err := pgx.Connect(CTX, config.Config.Database.URL())
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to database")
		return err
	}

	DB = conn
	Queries = repository.New(conn)

	return nil
}

func Close() {
	err := DB.Close(CTX)
	if err != nil {
		log.Error().Err(err).Msg("failed to close database")
		return
	}
}
