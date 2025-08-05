package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/maxpeterkaya/peanut/common"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"path"
	"peanut/internal/buildinfo"
	"strconv"
)

type ConfigStruct struct {
	Version       string        `toml:"version"`
	Admin         admin         `toml:"admin"`
	Common        commonStruct  `toml:"common"`
	Database      database      `toml:"database"`
	Github        github        `toml:"github"`
	Authorization authorization `toml:"authorization"`
}

type database struct {
	User string `toml:"user"`
	Pass string `toml:"pass"`
	Name string `toml:"name"`
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

func (d *database) URL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", d.User, d.Pass, d.Host, d.Port, d.Name)
}

type github struct {
	Token        string   `toml:"token"`
	Repositories []string `toml:"repositories"`
	Owner        string   `toml:"owner"`
}

type authorization struct {
	PrometheusToken  string `toml:"prometheus_token"`
	GithubToken      string `toml:"github_token"`
	HealthCheckToken string `toml:"health_check_token"`
}

type admin struct {
	User string `toml:"user"`
	Pass string `toml:"pass"`
}

type commonStruct struct {
	EncryptionKey string `toml:"encryption_key"`

	EnableDatabase bool `toml:"enable_database"`

	// Options if certain API endpoints should have authentication
	AuthPrometheus    bool `toml:"auth_prometheus"`
	AuthHealthChek    bool `toml:"auth_health_check"`
	AuthGithubMetrics bool `toml:"auth_github_metrics"`
}

var (
	Config      *ConfigStruct
	IsReady     bool
	IsContainer bool
)

func Init() error {
	IsReady = false
	fileName := "config.toml"
	IsContainer = getEnvAsBool("IS_CONTAINER", false)
	exists := common.Exists(configFolder(fileName))

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
				EncryptionKey:     getEnv("ENCRYPT_KEY", common.GenerateKey(32)),
				AuthPrometheus:    false,
				AuthHealthChek:    false,
				AuthGithubMetrics: false,

				EnableDatabase: false,
			},
			Github: github{
				Repositories: []string{""},
			},
			Authorization: authorization{
				PrometheusToken:  common.GeneratePassword(32),
				GithubToken:      common.GeneratePassword(32),
				HealthCheckToken: common.GeneratePassword(32),
			},
		}

		file, err := os.Create(configFolder(fileName))
		if err != nil {
			log.Error().Err(err).Msg("Error creating config.toml")
			return err
		}
		defer file.Close()

		err = toml.NewEncoder(file).Encode(Config)
		if err != nil {
			log.Error().Err(err).Msg("Error encoding config.toml")
			return err
		}

		log.Info().Msg("generated new config.toml")
		log.Info().Msg("quitting application for user configuration...")

		os.Exit(0)
	} else {
		file, err := os.ReadFile(configFolder(fileName))
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

	if Config.Version != fmt.Sprintf("%s-%s", buildinfo.Version, buildinfo.Commit) {
		log.Info().Msg("updating config.toml to current schema...")

		backup, err := os.Create(configFolder("backup-" + fileName))
		if err != nil {
			log.Error().Err(err).Msg("Error creating backup-" + fileName)
			return err
		}
		defer backup.Close()

		file, err := os.Open(configFolder(fileName))
		if err != nil {
			log.Error().Err(err).Msg("Error opening config.toml")
			return err
		}
		defer file.Close()

		_, err = io.Copy(backup, file)
		if err != nil {
			log.Error().Err(err).Msg("Error copying config.toml")
			return err
		}

		// So for some reason the toml encoder doesn't accept os.Open but it accepts os.Create file descriptor

		Config.Version = fmt.Sprintf("%s-%s", buildinfo.Version, buildinfo.Commit)

		// Check that tokens have been created
		if Config.Authorization.PrometheusToken == "" {
			Config.Authorization.PrometheusToken = common.GeneratePassword(32)
		}
		if Config.Authorization.GithubToken == "" {
			Config.Authorization.GithubToken = common.GeneratePassword(32)
		}
		if Config.Authorization.HealthCheckToken == "" {
			Config.Authorization.HealthCheckToken = common.GeneratePassword(32)
		}

		file, err = os.Create(configFolder(fileName))
		if err != nil {
			log.Error().Err(err).Msg("Error opening config.toml")
			return err
		}
		defer file.Close()

		err = toml.NewEncoder(file).Encode(Config)
		if err != nil {
			log.Error().Err(err).Msg("Error encoding config.toml")
			return err
		}

		log.Info().Msg("updated config.toml to latest schema.")
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

func configFolder(fileName string) string {
	if IsContainer {
		return path.Clean("/config/" + fileName)
	} else {
		return path.Clean(fileName)
	}
}
