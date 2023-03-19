package validator

import (
	"fmt"

	"github.com/dmitrymomot/oauth2-server/internal/utils"
	"github.com/gookit/validate"
	"github.com/pkg/errors"
)

func init() {
	// change global opts
	validate.Config(func(opt *validate.GlobalOption) {
		opt.FieldTag = "json"
		opt.StopOnError = false
		opt.SkipOnEmpty = true
		opt.UpdateSource = true
		opt.CheckZero = false
		opt.ErrKeyFmt = 1
	})

	// Add custom global validation rules
	validate.AddValidators(validate.M{
		"realEmail": func(val interface{}) bool {
			email, ok := val.(string)
			if !ok {
				return false
			}

			email, err := utils.SanitizeEmail(email)
			if err != nil {
				return false
			}

			if err := ValidateEmail(email); err != nil {
				return false
			}

			return true
		},
	})

	// Add global filters
	validate.AddFilters(validate.M{
		"sanitizeEmail": func(val interface{}) (string, error) {
			if email, ok := val.(string); ok {
				return utils.SanitizeEmail(email)
			}

			return fmt.Sprintf("%v", val), errors.New("invalid email address")
		},
	})

	// Add global messages
	validate.AddGlobalMessages(map[string]string{
		"realEmail":     "Email address is not real",
		"sanitizeEmail": "Invalid email address",
	})
}
