package validator

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/dmitrymomot/oauth2-server/internal/utils"
)

const (
	// Predefined image content types
	ImageContentTypeJpeg = "image/jpeg"
	ImageContentTypePng  = "image/png"
	ImageContentTypeGif  = "image/gif"
)

// Predefined errors
var (
	ErrFileSizeLt             = errors.New("file size is less than")
	ErrFileSizeGt             = errors.New("file size is greater than")
	ErrInvalidFileContentType = errors.New("invalid file content type")
)

type (
	// fileValidationFunc is a function interface for validating file size and content type
	fileValidationFunc func(*multipart.FileHeader) error
)

// DefaultFormMaxMemory is the default maximum amount of memory to use when parsing a multipart form.
var DefaultFormMaxMemory int64 = 32 << 20 // 32 MB

// ValidateFile validates file size and content type.
// It returns url.Values if there is an error otherwise it returns nil.
func ValidateFile(file *multipart.FileHeader, field string, validators ...fileValidationFunc) url.Values {
	errs := url.Values{}
	for _, validator := range validators {
		if err := validator(file); err != nil {
			errs.Add(field, utils.UcFirst(err.Error()))
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

// ValidateImage validates image file size and content type.
// It returns url.Values if there is an error otherwise it returns nil.
func ValidateImage(file *multipart.FileHeader, field, min, max string) url.Values {
	return ValidateFile(file, field,
		ValidateFileSize(min, max),
		ValidateFileContentType(ImageContentTypeJpeg, ImageContentTypePng),
	)
}

// ValidateFileFromRequest validates file size and content type from multipart form request.
// It returns url.Values if there is an error otherwise it returns nil.
func ValidateFileFromRequest(r *http.Request, field string, validators ...fileValidationFunc) url.Values {
	_, fileHeader, err := r.FormFile(field)
	if err != nil {
		return url.Values{field: {err.Error()}}
	}

	return ValidateFile(fileHeader, field, validators...)
}

// ValidateFileSize validates file size.
func ValidateFileSize(min, max string) fileValidationFunc {
	return func(file *multipart.FileHeader) error {
		if file.Size == 0 {
			return errors.New("file size is zero")
		}

		minSize, err := utils.ParseFileSize(min)
		if err != nil {
			return fmt.Errorf("invalid min file size: %s", err)
		}

		maxSize, err := utils.ParseFileSize(max)
		if err != nil {
			return fmt.Errorf("invalid max file size: %s", err)
		}

		if file.Size < minSize {
			return fmt.Errorf("%s %s", ErrFileSizeLt, min)
		}

		if file.Size > maxSize {
			return fmt.Errorf("%s %s", ErrFileSizeGt, max)
		}

		return nil
	}
}

// ValidateFileContentType validates file content type.
// It returns url.Values if there is an error otherwise it returns nil.
func ValidateFileContentType(contentTypes ...string) fileValidationFunc {
	return func(file *multipart.FileHeader) error {
		for _, contentType := range contentTypes {
			if file.Header.Get("Content-Type") == contentType {
				return nil
			}
		}

		return ErrInvalidFileContentType
	}
}
