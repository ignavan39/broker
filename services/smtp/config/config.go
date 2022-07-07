package config

import "os"

type Config struct {
	Email string 
	Password string
	Host string
	Address string
}

var config Config

func Init() {

	config = Config{
		Email: os.Getenv("USER_GMAIL"),
		Password: os.Getenv("USER_PASSWORD"),
		Host: "smtp.gmail.com",
		Address: "smtp.gmail.com:465",
	}
}

func GetConfig() Config {
	return config
}