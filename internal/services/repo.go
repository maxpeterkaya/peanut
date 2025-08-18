package services

import (
	"peanut/internal/database"
	"peanut/internal/repository"
	"strconv"

	"github.com/rs/zerolog/log"
)

func CreateRepo() {

}

func SearchRepo(search string) (*repository.Repository, error) {
	isNumber := false
	num, err := strconv.ParseUint(search, 10, 64)
	if err == nil {
		isNumber = true
	}

	if isNumber {
		repo, err := database.Queries.GetRepository(database.CTX, int32(num))
		if err != nil {
			log.Error().Err(err).Msg("Failed to get repository")
			return nil, err
		}
		return &repo, nil
	} else {
		repo, err := database.Queries.SearchRepository(database.CTX, search)
		if err != nil {
			log.Error().Err(err).Msg("Failed to search repository")
			return nil, err
		}
		return &repo, nil
	}
}
