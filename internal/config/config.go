package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

const (
	defaultPort       = "8080"
	defaultDBHost     = "localhost"
	defaultDBPort     = "5432"
	defaultDBUser     = "postgres"
	defaultDBPassword = "password"
	defaultDBName     = "subscriptions"
	defaultLogLevel   = "info"
)

type Config struct {
	Server   server
	database database
	Logger   logger
}

type server struct {
	port string `env:"PORT"`
}

type database struct {
	dbHost     string `env:"DB_HOST"`
	dbPort     string `env:"DB_PORT"`
	dbUser     string `env:"DB_USER"`
	dbPassword string `env:"DB_PASSWORD"`
	dbName     string `env:"DB_NAME"`
}

type logger struct {
	logLevel string `env:"LOG_LEVEL"`
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c Config) GetDatabaseURL() string {
	databaseURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.database.dbHost,
		c.database.dbPort,
		c.database.dbUser,
		c.database.dbPassword,
		c.database.dbName)

	return databaseURL
}
