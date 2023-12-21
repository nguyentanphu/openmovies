package main

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/v1/healthcheck", app.healthcheckHandler).Methods(http.MethodGet)
	router.HandleFunc("/v1/movies", app.createMovie).Methods(http.MethodPost)
	router.HandleFunc("/v1/movies/{id:[0-9]+}", app.getMovie).Methods(http.MethodGet)

	router.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:4000/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	return router
}
