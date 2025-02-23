package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Service  Service
	Auth     Auth
	QR       QR
	Metrics  Metrics
	Platform Platform
}

type Service struct {
	Port string `env:"GATEWAY_SERVICE_PORT"`
	Name string `env:"GATEWAY_SERVICE_NAME"`
}

type Auth struct {
	Host string `env:"AUTH_SERVICE_HOST"`
	Port string `env:"AUTH_SERVICE_PORT"`
}

type QR struct {
	Host string `env:"QR_SERVICE_HOST"`
	Port string `env:"QR_SERVICE_PORT"`
}

type Metrics struct {
	Host string `env:"GRAFANA_HOST"`
	Port int    `env:"GRAFANA_PORT"`
}

type Platform struct {
	SigningKey string `env:"SECURITY_AUTH_SIGNING_KEY"`
	Env        string `env:"ENV"`
}

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)

	if err != nil {
		log.Fatalf("Can not read env variables: %s", err)
	}

	return cfg
}
