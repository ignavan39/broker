package config

import (
	"broker/pkg/pg"
	"fmt"
	"os"
	"strconv"
	"time"
)

type JWTConfig struct {
	HashSalt       string        `env:"HASH_SALT" envDefault:"super_secret"`
	SigningKey     string        `env:"SIGNING_KEY" envDefault:"signing_key"`
	ExpireDuration time.Duration `env:"EXPIRE_DURATION" envDefault:"24h"`
}

type AMQPConfig struct {
	Host             string `env:"AMQP_HOST" envDefault:"broker.rabbitmq.loc"`
	Port             int    `env:"AMQP_PORT" envDefault:"5672"`
	User             string `env:"AMQP_USER" envDefault:"user"`
	Pass             string `env:"AMQP_PASS" envDefault:"pass"`
	ExternalUser     string `env:"AMQP_EXTERNAL_USER" envDefault:"user"`
	ExternalPassword string `env:"AMQP_EXTERNAL_PASS" envDefault:"pass"`
	Vhost            string `env:"AMQP_VHOST" envDefault:"/"`
	QueueHashSalt    string `env:"QUEUE_HASH_SALT" envDefault:"super_secret"`
}

type Config struct {
	JWT      JWTConfig
	Database pg.Config
	AMQP     AMQPConfig
}

var config = Config{}

func Init() error {
	dbPort, err := strconv.ParseInt(os.Getenv("DATABASE_PORT"), 10, 16)
	if err != nil {
		return fmt.Errorf("error for parsing DATABASE_PORT :%s", err)
	}

	pgConf := pg.Config{
		Password: os.Getenv("DATABASE_PASS"),
		Host:     os.Getenv("DATABASE_HOST"),
		User:     os.Getenv("DATABASE_USER"),
		Port:     uint16(dbPort),
		DB:       os.Getenv("DATABASE_NAME"),
	}

	amqpPort, err := strconv.ParseInt(os.Getenv("AMQP_PORT"), 10, 16)
	if err != nil {
		return fmt.Errorf("error for parsing AMQP_PORT :%s", err)
	}
	amqpConf := AMQPConfig{
		Port:             int(amqpPort),
		Host:             os.Getenv("AMQP_HOST"),
		User:             os.Getenv("AMQP_USER"),
		Pass:             os.Getenv("AMQP_PASS"),
		Vhost:            os.Getenv("AMQP_VHOST"),
		ExternalUser:     os.Getenv("AMQP_EXTERNAL_USER"),
		ExternalPassword: os.Getenv("AMQP_EXTERNAL_PASS"),
		QueueHashSalt:    os.Getenv("QUEUE_HASH_SALT"),
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
		Database: pgConf,
		JWT:      jwt,
		AMQP:     amqpConf,
	}
	return nil
}

func GetConfig() Config {
	return config
}
