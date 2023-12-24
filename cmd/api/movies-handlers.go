package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"openmovies/internal/data"
	"strconv"
)

type MovieDto struct {
	Title   string       `json:"title" validate:"required,max=500"`
	Year    int32        `json:"year" validate:"required,min=1888"`
	Runtime data.Runtime `json:"runtime" validate:"gt=0"`
	Genres  []string     `json:"genres" validate:"required,min=1,max=5,unique"`
}

func (app *application) postMovieHandler(w http.ResponseWriter, r *http.Request) {
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
	movie := &data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}
	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = app.writeJson(w, http.StatusOK, envelop{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	movie, err := app.models.Movies.GetById(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)

		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJson(w, http.StatusOK, envelop{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) putMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.GetById(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)

		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	var input MovieDto
	err = app.decodeJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	apiErr := app.validateInput(input)
	if apiErr != nil {
		app.fieldValidationResponse(w, r, apiErr)
		return
	}
	movie.Title = input.Title
	movie.Year = input.Year
	movie.Runtime = input.Runtime
	movie.Genres = input.Genres

	err = app.models.Movies.Update(movie)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJson(w, http.StatusOK, envelop{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJson(
		w,
		http.StatusOK,
		envelop{"message": fmt.Sprintf("movie with id: %d was deleted", id)},
		nil,
	)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

type MoviePartialDto struct {
	Title   *string       `json:"title" validate:"omitnil,max=500"`
	Year    *int32        `json:"year" validate:"omitnil,min=1888"`
	Runtime *data.Runtime `json:"runtime" validate:"omitnil,gt=0"`
	Genres  []string      `json:"genres" validate:"omitnil,min=1,max=5,unique"`
}

func (app *application) patchMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(mux.Vars(r)["id"], 10, 64)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.GetById(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)

		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input MoviePartialDto
	err = app.decodeJson(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		movie.Title = *input.Title
	}
	if input.Year != nil {
		movie.Year = *input.Year
	}
	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}
	if input.Genres != nil {
		movie.Genres = input.Genres
	}

	apiErr := app.validateInput(input)
	if apiErr != nil {
		app.fieldValidationResponse(w, r, apiErr)
		return
	}

	err = app.models.Movies.Update(movie)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJson(w, http.StatusOK, envelop{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getMovies(w http.ResponseWriter, r *http.Request) {
	input := data.NewMovieFilters()
	err := app.schemaDecoder.Decode(&input, r.URL.Query())
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	apiErr := app.validateInput(input)
	if apiErr != nil {
		app.fieldValidationResponse(w, r, apiErr)
		return
	}
	movies, err := app.models.Movies.Get(input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJson(w, http.StatusOK, envelop{"movies": movies}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
