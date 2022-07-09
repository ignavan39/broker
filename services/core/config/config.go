package config

import (
	"broker/pkg/pg"
	"fmt"
	"os"
	"strconv"
	"time"
)

type JWTConfig struct {
	HashSalt              string        `env:"HASH_SALT" envDefault:"super_secret"`
	SigningKey            string        `env:"SIGNING_KEY" envDefault:"signing_key"`
	AccessExpireDuration  time.Duration `env:"ACCESS_EXPIRE_DURATION" envDefault:"30m"`
	RefreshExpireDuration time.Duration `env:"REFRESH_EXPIRE_DURATION" envDefault:"168h"`
}

type RedisConfig struct {
	Host     string `env:"REDIS_HOST"`
	Port     int    `env:"REDIS_PORT"`
	DB       int    `env:"REDIS_DB"`
	Password string `env:"REDIS_PASSWORD"`
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

type MailgunConfig struct {
	PrivateKey string `env:"MAILGUN_API_KEY"`
	Domain     string `env:"MAILGUN_DOMAIN"`
	PublicKey  string `env:"MAILGUN_PUBLIC_KEY"`
	Sender     string `env:"MAILGUN_SENDER"`
}

type FrontendConfig struct {
	Host string `env:"FRONTEND_HOST" envDefault:"localhost:3000"`
}

type InvitationConfig struct {
	InvitationHashSalt       string        `env:"INVITATION_HASH_SALT" envDefault:"puper_secret"`
	InvitationExpireDuration time.Duration `env:"INVITATION_EXPIRE_DURATION" envDefault:"5s"`
}

type Config struct {
	JWT           JWTConfig
	Database      pg.Config
	MailgunConfig MailgunConfig
	AMQP          AMQPConfig
	Redis         RedisConfig
	Frontend      FrontendConfig
	Invitation    InvitationConfig
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

	accessExpireDurationRaw := os.Getenv("ACCESS_EXPIRE_DURATION")
	accessExpireDuration, err := time.ParseDuration(accessExpireDurationRaw)

	if err != nil {
		return fmt.Errorf("error for parsing ACCESS_EXPIRE_DURATION :%s", err)
	}

	refreshExpireDurationRaw := os.Getenv("REFRESH_EXPIRE_DURATION")
	refreshExpireDuration, err := time.ParseDuration(refreshExpireDurationRaw)

	if err != nil {
		return fmt.Errorf("error for parsing REFRESH_EXPIRE_DURATION :%s", err)
	}

	jwt := JWTConfig{
		HashSalt:              os.Getenv("HASH_SALT"),
		SigningKey:            os.Getenv("SIGNING_KEY"),
		AccessExpireDuration:  accessExpireDuration,
		RefreshExpireDuration: refreshExpireDuration,
	}

	mailgunConfig := MailgunConfig{
		Domain:     os.Getenv("MAILGUN_DOMAIN"),
		PrivateKey: os.Getenv("MAILGUN_API_KEY"),
		PublicKey:  os.Getenv("MAILGUN_PUBLIC_KEY"),
		Sender:     os.Getenv("MAILGUN_SENDER"),
	}

	redisPort, err := strconv.ParseInt(os.Getenv("REDIS_PORT"), 10, 16)
	if err != nil {
		return fmt.Errorf("error for parsing REDIS_PORT :%s", err)
	}

	redisDB, err := strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 16)
	if err != nil {
		return fmt.Errorf("error for parsing REDIS_DB :%s", err)
	}

	redis := RedisConfig{
		Host:     os.Getenv("REDIS_HOST"),
		Port:     int(redisPort),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       int(redisDB),
	}

	frontend := FrontendConfig{
		Host: os.Getenv("FRONTEND_HOST"),
	}

	invitationExpireDurationRaw := os.Getenv("INVITATION_EXPIRE_DURATION")
	invitationExpireDuration, err := time.ParseDuration(invitationExpireDurationRaw)
	if err != nil {
		return fmt.Errorf("error for parsing INVITATION_EXPIRE_DURATION :%s", err)
	}

	invitation := InvitationConfig{
		InvitationHashSalt:       os.Getenv("INVITATION_HASH_SALT"),
		InvitationExpireDuration: invitationExpireDuration,
	}

	config = Config{
		Database:      pgConf,
		JWT:           jwt,
		MailgunConfig: mailgunConfig,
		AMQP:          amqpConf,
		Redis:         redis,
		Frontend:      frontend,
		Invitation:    invitation,
	}
	return nil
}

func GetConfig() Config {
	return config
}
