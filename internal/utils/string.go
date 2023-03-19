package utils

import (
	"strings"

	"github.com/labstack/gommon/bytes"
	"github.com/pkg/errors"
)

// Trim string between two substrings and return the string without it and substrings.
func TrimStringBetween(str, start, end string) string {
	indx1 := strings.Index(str, start)
	indx2 := strings.Index(str, end)
	if indx1 == -1 || indx2 == -1 {
		return strings.TrimSpace(str)
	}
	return strings.TrimSpace(str[:indx1] + str[indx2+len(end):])
}

// TrimRightZeros trims trailing zeros from string.
func TrimRightZeros(str string) string {
	return strings.TrimRight(str, "0")
}

// Parse file size from string to int64
func ParseFileSize(s string) (int64, error) {
	size, err := bytes.Parse(s)
	if err != nil {
		return 0, errors.Wrap(err, "failed to parse file size")
	}

	return size, nil
}

// UcFirst capitalizes first letter of a string
func UcFirst(s string) string {
	if len(s) == 0 {
		return s
	}

	return strings.ToUpper(s[:1]) + s[1:]
}
