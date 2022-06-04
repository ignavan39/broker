package config

import (
	"broker/pkg/pg"
	"fmt"
	"os"
	"strconv"
)


type Config struct {
	Database    pg.Config
	Port        int `env:"PORT"`}

var config = Config{}

func Init() (error) {
	dbPort, err := strconv.ParseInt(os.Getenv("DATABASE_PORT"), 10, 16)
	if err != nil {
		return fmt.Errorf("error for parsing DATABASE_PORT :%s",err)
	}
	
	port,err := strconv.ParseInt(os.Getenv("PORT"), 10, 16)
	if err != nil {
		return fmt.Errorf("error for parsing PORT :%s",err)
	}

	pgCong := pg.Config{
		Password: os.Getenv("DATABASE_PASS"),
		Host:     os.Getenv("DATABASE_HOST"),
		User:     os.Getenv("DATABASE_USER"),
		Port:     uint16(dbPort),
		DB:       os.Getenv("DATABASE_NAME"),
	}
	
	config = Config{
		Database: pgCong,
		Port: int(port),
	}
	return nil
}

func GetConfig() Config {
	return config
}
