package config

import (
	"github.com/BurntSushi/toml"
	"github.com/maxpeterkaya/peanut/common"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

type ConfigStruct struct {
	Admin    admin        `toml:"admin"`
	Common   commonStruct `toml:"common"`
	Database database     `toml:"database"`
	Github   github       `toml:"github"`
}

type database struct {
	User string `toml:"user"`
	Pass string `toml:"pass"`
	Name string `toml:"name"`
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type github struct {
	Token      string `toml:"token"`
	Repository string `toml:"repository"`
	RepoOwner  string `toml:"repo_owner"`
}

type admin struct {
	User string `toml:"user"`
	Pass string `toml:"pass"`
}

type commonStruct struct {
	EncryptionKey string `toml:"encryption_key"`
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
			Database: database{
				Host: getEnv("DB_HOST", "localhost"),
				User: getEnv("DB_USER", "postgres"),
				Pass: getEnv("DB_PASS", "postgres"),
				Port: getEnvAsInt("DB_PORT", 5432),
				Name: getEnv("DB_NAME", "peanut"),
			},
			Admin: admin{
				User: getEnv("ADMIN_USER", common.GenerateUsername()),
				Pass: getEnv("ADMIN_PASS", common.GeneratePassword(10)),
			},
			Common: commonStruct{
				EncryptionKey: getEnv("ENCRYPT_KEY", common.GenerateKey(32)),
			},
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
			log.Error().Err(err).Msg("Error opening econfig.toml")
			return err
		}

		if _, err := toml.Decode(string(file), &Config); err != nil {
			log.Error().Err(err)
			return err
		}

		log.Info().Msg("loaded econfig.toml")
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
