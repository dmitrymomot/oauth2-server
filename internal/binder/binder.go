package binder

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/form/v4"
)

// BindJSON binds the json request body to the given interface.
func BindJSON[T any](r *http.Request, i T) error {
	if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
		return fmt.Errorf("failed to decode json: %w", err)
	}
	return nil
}

// BindForm binds the form request body to the given interface.
func BindForm[T any](r *http.Request, i T) error {
	if r.Form == nil {
		if err := r.ParseForm(); err != nil {
			return fmt.Errorf("failed to parse form: %w", err)
		}
	}

	decoder := form.NewDecoder()
	decoder.SetTagName("json")
	if err := decoder.Decode(&i, r.Form); err != nil {
		return fmt.Errorf("failed to decode form: %w", err)
	}

	return nil
}

// BindQuery binds the query string to the given interface.
func BindQuery[T any](r *http.Request, i T) error {
	decoder := form.NewDecoder()
	decoder.SetTagName("json")
	if err := decoder.Decode(&i, r.URL.Query()); err != nil {
		return fmt.Errorf("failed to decode query: %w", err)
	}

	return nil
}

// Bind binds the request body to the given interface.
// It will try to bind the request data depending on the request content type.
func Bind[T any](r *http.Request, i T) error {
	switch r.Header.Get("Content-Type") {
	case "application/json":
		return BindJSON(r, i)
	case "application/x-www-form-urlencoded", "multipart/form-data":
		return BindForm(r, i)
	default:
		return BindQuery(r, i)
	}
}
