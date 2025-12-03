package validate

import (
    "regexp"

    "github.com/go-playground/validator/v10"
)

// Example: simple mainland China mobile pattern.
// Adjust according to your real requirements.
var mobileRegexp = regexp.MustCompile(`^1[3-9]\d{9}$`)

func init() {
    // Register tag `mobile` for phone number validation.
    MustRegister("mobile", func(fl validator.FieldLevel) bool {
        v := fl.Field().String()
        if v == "" {
            // Let `required` handle emptiness if needed.
            return true
        }
        return mobileRegexp.MatchString(v)
    })
}
