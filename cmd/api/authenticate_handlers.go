package main

import (
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strconv"
	"time"
)

func (app *application) authenticateHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}

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

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCredentialResponse(w, r)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": strconv.FormatInt(user.ID, 10),
		"iat": jwt.NewNumericDate(time.Now()),
		"nbf": jwt.NewNumericDate(time.Now()),
		"exp": jwt.NewNumericDate(time.Now().Add(time.Hour)),
		"iss": "github.com/nguyentanphu/openmovies",
		"aud": "github.com/nguyentanphu/openmovies",
	})

	tokenString, err := token.SignedString([]byte(app.config.jwtSecret))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJson(w, http.StatusOK, envelop{
		"authentication-token": tokenString,
	}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
