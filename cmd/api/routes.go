package main

import (
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/v1/healthcheck", app.requireActivatedUser(app.healthcheckHandler)).Methods(http.MethodGet)

	router.HandleFunc("/v1/movies", app.requirePermission("movies:read", app.getMovies)).Methods(http.MethodGet)
	router.HandleFunc("/v1/movies", app.requirePermission("movies:write", app.postMovieHandler)).Methods(http.MethodPost)
	router.HandleFunc("/v1/movies/{id:[0-9]+}", app.requirePermission("movies:write", app.getMovieHandler)).Methods(http.MethodGet)
	router.HandleFunc("/v1/movies/{id:[0-9]+}", app.requirePermission("movies:write", app.putMovieHandler)).Methods(http.MethodPut)
	router.HandleFunc("/v1/movies/{id:[0-9]+}", app.requirePermission("movies:write", app.patchMovieHandler)).Methods(http.MethodPatch)
	router.HandleFunc("/v1/movies/{id:[0-9]+}", app.requirePermission("movies:write", app.deleteMovieHandler)).Methods(http.MethodDelete)

	router.HandleFunc("/v1/users", app.registerUserHandler).Methods(http.MethodPost)
	router.HandleFunc("/v1/users/activate", app.activateUserHandler).Methods(http.MethodPut)
	router.HandleFunc("/v1/users/auth", app.authenticateHandler).Methods(http.MethodPut)

	router.MethodNotAllowedHandler = http.HandlerFunc(app.methodNotAllowedResponse)
	router.NotFoundHandler = http.HandlerFunc(app.notFoundResponse)

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:4000/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	middlewares := alice.New(app.rateLimit, app.authenticate, app.recoverPanic)

	return middlewares.Then(router)
}
