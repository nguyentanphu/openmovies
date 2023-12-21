package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"openmovies/internal/data"
	"strconv"
)

type MovieDto struct {
	Title   string       `json:"title"`
	Year    int32        `json:"year"`
	Runtime data.Runtime `json:"runtime"`
	Genres  []string     `json:"genres"`
}

func (app *application) createMovie(w http.ResponseWriter, r *http.Request) {
	var input MovieDto
	err := app.decodeJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
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
