package conf

import (
	"github.com/happyreturns/gohelpers/config"
)

type Conf struct {
	config.DefaultConfig
	Port       string `env:"PORT"`
	TMDBUrl    string `env:"TMDB_URL"`
	TMDBApiKey string `env:"TMDB_API_KEY"`
}
