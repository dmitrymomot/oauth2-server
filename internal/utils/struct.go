package utils

import (
	"fmt"
	"net/url"

	"github.com/google/go-querystring/query"
)

// StructToUrlValues converts a struct to a url.Values object.
func StructToUrlValues(v interface{}) (url.Values, error) {
	if v == nil {
		return nil, fmt.Errorf("struct is nil")
	}

	if _, ok := v.(url.Values); ok {
		return v.(url.Values), nil
	}

	uv, err := query.Values(v)
	if err != nil {
		return nil, fmt.Errorf("failed to convert struct to url values: %w", err)
	}

	return uv, nil
}
