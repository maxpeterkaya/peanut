package config

import (
	"github.com/BurntSushi/toml"
	"github.com/maxpeterkaya/peanut/common"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

type ConfigStruct struct {
	DbUser string `toml:"db_user"`
	DbPass string `toml:"db_pass"`
	DbName string `toml:"db_name"`
	DbHost string `toml:"db_host"`
	DbPort int    `toml:"db_port"`

	AdminUser string `toml:"admin_user"`
	AdminPass string `toml:"admin_pass"`

	EncryptionKey string `toml:"encryption_key"`

	GHToken      string `toml:"gh_token"`
	GHRepository string `toml:"gh_repository"`
	GHRepoOwner  string `toml:"gh_repo_owner"`
}

var (
	Config  *ConfigStruct
	IsReady bool
)

func Init() error {
	IsReady = false
	fileName := "config.toml"
	exists := common.Exists(fileName)

	if !exists {
		Config = &ConfigStruct{
			DbHost: getEnv("DB_HOST", "localhost"),
			DbUser: getEnv("DB_USER", "postgres"),
			DbPass: getEnv("DB_PASS", "postgres"),
			DbPort: getEnvAsInt("DB_PORT", 5432),
			DbName: getEnv("DB_NAME", "peanut"),

			AdminUser: getEnv("ADMIN_USER", common.GenerateUsername()),
			AdminPass: getEnv("ADMIN_PASS", common.GeneratePassword(10)),

			EncryptionKey: getEnv("ENCRYPT_KEY", common.GenerateKey(32)),
		}

		file, err := os.Create(fileName)
		if err != nil {
			log.Error().Err(err).Msg("Error creating config.yml")
			return err
		}

		err = toml.NewEncoder(file).Encode(Config)
		if err != nil {
			log.Error().Err(err).Msg("Error encoding config.yml")
			return err
		}

		log.Info().Msg("generated new config.yml")
		IsReady = true
	} else {
		file, err := os.ReadFile(fileName)
		if err != nil {
			log.Error().Err(err).Msg("Error opening config.toml")
			return err
		}

		if _, err := toml.Decode(string(file), &Config); err != nil {
			log.Error().Err(err)
			return err
		}

		log.Info().Msg("loaded config.toml")
		IsReady = true
	}

	return nil
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}

func getEnvAsBool(name string, defaultVal bool) bool {
	valueStr := getEnv(name, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}

	return defaultVal
}
