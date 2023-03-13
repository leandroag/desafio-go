package app

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Postgres PostgresConfig
}

type PostgresConfig struct {
	Host        string `env:"DB_HOST" default:"localhost"`
	Port        string `env:"DB_PORT" default:"5432"`
	Database    string `env:"DB_DATABASE" default:"postgres"`
	User        string `env:"DB_USER" default:"postgres"`
	Password    string `env:"DB_PASS" default:"postgres"`
	PoolMinSize string `env:"DB_POOL_MIN_SIZE" default:"2"`
	PoolMaxSize string `env:"DB_POOL_MAX_SIZE" default:"10"`
	SSLMode     string `env:"DB_SSL_MODE" default:"disable"`
}

func ReadConfigFromEnv() *Config {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("error reading env")
	}

	return &cfg
}

func ReadConfigFromFile(filename string) *Config {
	var cfg Config
	err := cleanenv.ReadConfig(filename, &cfg)
	if err != nil {
		log.Fatal().Stack().Err(err).Msg("error reading file")
	}

	return &cfg
}

func ReadConfig(filename string) *Config {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Warn().Msgf("File not found %s", filename)
		return ReadConfigFromEnv()
	}

	return ReadConfigFromFile(filename)
}

func (pg PostgresConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		pg.Host,
		pg.Port,
		pg.Database,
		pg.User,
		pg.Password,
		pg.SSLMode,
	)
}

func (pg PostgresConfig) URL() string {
	if pg.SSLMode == "" {
		pg.SSLMode = "disable"
	}

	connectString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		pg.User, pg.Password, pg.Host, pg.Port, pg.Database, pg.SSLMode)

	return connectString
}
