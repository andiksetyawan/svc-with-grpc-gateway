package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog/log"
)

type configs struct {
	ServiceName           string `env:"SERVICE_NAME" yaml:"service_name" env-default:"svc-with-grpc-gateway"`
	Address               string `env:"ADDRESS" yaml:"address" env-default:":8080"`
	OtlpCollectorUrl      string `env:"OTLP_COLLECTOR_URL" yaml:"otlp_collector_url" env-default:"localhost:4317"`
	InsecureOtlpCollector string `env:"INSECURE_OTLP_COLLECTOR" yaml:"insecure_otlp_collector" env-default:"true"`
}

var App configs

func Init() {
	err := cleanenv.ReadConfig(".env", &App)
	if err != nil {
		log.Debug().Msg("failed to read .env file, setting config from environment variables.")
		cleanenv.ReadEnv(&App)
	}
}
