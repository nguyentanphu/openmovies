package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"openmovies/internal/data"
	"strconv"
)

type MovieDto struct {
	Title   string       `json:"title" validate:"max=500"`
	Year    int32        `json:"year" validate:"min=1888"`
	Runtime data.Runtime `json:"runtime" validate:"gt=0"`
	Genres  []string     `json:"genres" validate:"required,min=1,max=5,unique"`
}

func (app *application) createMovie(w http.ResponseWriter, r *http.Request) {
	var input MovieDto
	err := app.decodeJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	apiErr := app.validateInput(input)
	if apiErr != nil {
		app.fieldValidationResponse(w, r, apiErr)
		return
	}
	fmt.Fprintf(w, "%+v", input)
}

func (app *application) getMovie(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.NotFound(w, r)
		return
	}
	movie := data.Movie{
		ID:      int64(id),
		Title:   "Test best movie",
		Runtime: 102,
		Genres:  []string{"Thriller", "Sci-fi"},
		Version: 1,
		Year:    2001,
	}

	err = app.writeJson(w, http.StatusOK, envelop{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
