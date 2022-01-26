package main

import (
	api2 "github.com/happyreturns/api-testing/api"
	conf2 "github.com/happyreturns/api-testing/conf"
	"github.com/happyreturns/gohelpers/config"
	"github.com/happyreturns/gohelpers/log"
)

func main() {
	conf := &conf2.Conf{
		DefaultConfig: config.DefaultConfig{},
		Port:          "7777",
		TMDBApiKey:    "{REDACTED}",
		TMDBUrl:       "https://api.themoviedb.org",
	}
	logger := log.NewLogger(conf.App, conf.Environment)

	config.FillConfig(conf)

	api := api2.NewApi(conf, logger)
	api.Initialize()

	api.Run()
}
