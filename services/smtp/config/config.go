package config

import "os"

type Config struct {
	Email    string `env:"SENDER_GMAIL" envDefault:"infbrokerinfo@gmail.com"`
	Password string `env:"SENDER_GMAIL_PASSWORD" envDefault:"pyuvwbfwumkkqrfd"`
	Host     string
	Address  string
}

var config Config

func Init() {

	config = Config{
		Email:    os.Getenv("SENDER_GMAIL"),
		Password: os.Getenv("SENDER_PASSWORD_GMAIL"),
		Host:     "smtp.gmail.com:",
		Address:  "smtp.gmail.com:465",
	}
}

func GetConfig() Config {
	return config
}
