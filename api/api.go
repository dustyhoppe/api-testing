package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/happyreturns/api-testing/api/handlers"
	"github.com/happyreturns/api-testing/clients"
	"github.com/happyreturns/api-testing/conf"
	"github.com/happyreturns/gohelpers/log"
	"github.com/rs/cors"
	"net/http"
	"time"
)

type Api struct {
	Config        *conf.Conf
	Logger        *log.Logger
	Handler       http.Handler
	MovieApi      *clients.MovieApiClient
	MovieHandlers *handlers.MovieHandlers
}

func NewApi(conf *conf.Conf, logger *log.Logger) *Api {

	movieApi := clients.NewMovieApiClient(conf, logger)
	movieHandlers := handlers.NewMovieHandlers(movieApi, logger)

	return &Api{
		Config:        conf,
		Logger:        logger,
		MovieApi:      movieApi,
		MovieHandlers: movieHandlers,
	}
}

func (api *Api) Initialize() {
	router := mux.NewRouter()

	router.HandleFunc("/v1/movies/search", api.MovieHandlers.SearchMoviesHandler).Methods("GET")
	router.HandleFunc("/v1/sessions", api.MovieHandlers.CreateGuestSession).Methods("POST")

	c := cors.AllowAll()
	api.Handler = c.Handler(router)
}

func (api *Api) Run() {

	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", api.Config.Port),
		Handler:      api.Handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	api.Logger.Fatal(server.ListenAndServe())
}
