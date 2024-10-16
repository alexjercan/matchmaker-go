package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Database struct {
		Host     string `env:"DB_HOST"`
		Port     int    `env:"DB_PORT"`
		User     string `env:"DB_USER"`
		Password string `env:"DB_PASSWORD"`
		Database string `env:"DB_DATABASE"`
		Insecure bool   `env:"DB_INSECURE"`
	}
}

func LoadConfig() Config {
	cfg := Config{}

	env := os.Getenv("ENV")
	if "" == env {
		env = "local"
	}

	err := godotenv.Load(".env." + env)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatal("Error read env variables")
	}

	return cfg
}

