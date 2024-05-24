package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Env string `env:"ENV"`
	URL string `env:"POSTGRES_URL"`
}

func NewConfig() (*Config, error) {
	cfg := new(Config)

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
