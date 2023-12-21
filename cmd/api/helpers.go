package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type envelop map[string]interface{}

func (app *application) writeJson(w http.ResponseWriter, status int, data envelop, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func (app *application) decodeJson(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(data)
	var invalidUnmarshal *json.InvalidUnmarshalError

	if err != nil {
		if errors.As(err, &invalidUnmarshal) {
			panic(err)
		}
		return errors.New("error parsing JSON")
	}

	err = decoder.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}
	return nil
}
