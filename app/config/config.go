package config

import (
	"broker/pkg/pg"
	"fmt"
	"os"
	"strconv"
	"time"
)

type JWTConfig struct {
	HashSalt       string        `env:"HASH_SALT"`
	SigningKey     string        `env:"SIGNING_KEY"`
	ExpireDuration time.Duration `env:"EXPIRE_DURATION"`
}

type Config struct {
	JWT      JWTConfig
	Database pg.Config
}

var config = Config{}

func Init() error {
	dbPort, err := strconv.ParseInt(os.Getenv("DATABASE_PORT"), 10, 16)
	if err != nil {
		return fmt.Errorf("error for parsing DATABASE_PORT :%s", err)
	}

	pgCong := pg.Config{
		Password: os.Getenv("DATABASE_PASS"),
		Host:     os.Getenv("DATABASE_HOST"),
		User:     os.Getenv("DATABASE_USER"),
		Port:     uint16(dbPort),
		DB:       os.Getenv("DATABASE_NAME"),
	}

	expireDurationRaw := os.Getenv("EXPIRE_DURATION")
	expireDuration, err := time.ParseDuration(expireDurationRaw)
	if err != nil {
		return fmt.Errorf("error for parsing EXPIRE_DURATION :%s", err)
	}

	jwt := JWTConfig{
		HashSalt:       os.Getenv("HASH_SALT"),
		SigningKey:     os.Getenv("SIGNING_KEY"),
		ExpireDuration: expireDuration,
	}

	config = Config{
		Database: pgCong,
		JWT:      jwt,
	}
	return nil
}

func GetConfig() Config {
	return config
}
