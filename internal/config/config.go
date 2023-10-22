package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	PostgresAuth `yaml:"postgres_auth"`
	NatsAuth     `yaml:"nats_auth"`
}

type PostgresAuth struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password"`
	DbName   string `yaml:"dbname" env-required:"true"`
}

type NatsAuth struct {
	ChannelName   string `yaml:"channel_name" env-default:"channel"`
	StanClusterId string `yaml:"stanClusterId" env-default:"test-cluster"`
	ClientId      string `yaml:"clientId" env-required:"true"`
}

func MustLoad() *Config {
	configPath := "/Users/sazhinsema/GolandProjects/WB_Tasks/WB_L0/config/config.yaml"

	var config Config

	err := cleanenv.ReadConfig(configPath, &config)
	if err != nil {
		log.Fatal("config problems")
	}

	return &config
}
