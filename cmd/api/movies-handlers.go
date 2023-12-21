package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"openmovies/internal/data"
	"strconv"
)

func (app *application) createMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "create movie")
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
