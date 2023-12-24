package main

import (
	"net/url"
	"strconv"
	"strings"
)

func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	v := qs.Get(key)
	if v == "" {
		return defaultValue
	}
	return v
}

func (app *application) readCommaSeparatedValues(qs url.Values, key string, defaultValue []string) []string {
	v := qs.Get(key)
	if v == "" {
		return defaultValue
	}

	return strings.Split(v, ",")
}

func (app *application) readInt(qs url.Values, key string, defaultValue int) (int, error) {
	v := qs.Get(key)
	if v == "" {
		return defaultValue, nil
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return 0, err
	}

	return i, nil
}
