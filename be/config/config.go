package config

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	PORT string `env:"APP_PORT" envDefault:":7000"`
	JWT  JWT
	DB   DB
}

type JWT struct {
	Secret string `env:"JWT_SECRET" envDefault:"super-secret"`
	Expiry int    `env:"JWT_EXPIRY" envDefault:"3600"`
}

type DB struct {
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     string `env:"DB_PORT" envDefault:"5432"`
	User     string `env:"DB_USER" envDefault:"postgres"`
	Password string `env:"DB_PASSWORD" envDefault:"postgres"`
	Database string `env:"DB_NAME" envDefault:"postgres"`
}

func New(file string) (*Config, error) {
	if err := godotenv.Load(file); err != nil {
		log.Printf("unable to load .env file: %e", err)
		return nil, err
	}
	config := &Config{}

	if err := env.Parse(config); err != nil {
		log.Printf("Failed to parse config: %v", err)
		return nil, err
	}

	return config, nil
}
