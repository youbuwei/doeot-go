package validate

import (
    "github.com/go-playground/validator/v10"
    "github.com/youbuwei/doeot-go/pkg/errs"
)

// v is a global validator instance.
var v = validator.New()

// Custom can be implemented by request DTOs to provide extra validation logic.
type Custom interface {
    Validate() error
}

// Struct validates a struct using `validate` tags + optional Custom.Validate().
// On failure it returns an *errs.Error with CodeBadRequest.
func Struct(s any) error {
    // 1) tag-based rules
    if err := v.Struct(s); err != nil {
        if verrs, ok := err.(validator.ValidationErrors); ok && len(verrs) > 0 {
            return errs.BadRequest(verrs[0].Error())
        }
        return errs.BadRequest(err.Error())
    }

    // 2) per-request custom rules
    if c, ok := s.(Custom); ok {
        if err := c.Validate(); err != nil {
            if e, ok := err.(*errs.Error); ok {
                return e
            }
            return errs.BadRequest(err.Error())
        }
    }

    return nil
}

// Register registers a custom tag validator, e.g. "mobile", "username".
func Register(tag string, fn validator.Func) error {
    return v.RegisterValidation(tag, fn)
}

// MustRegister registers a custom tag validator and panics on error.
func MustRegister(tag string, fn validator.Func) {
    if err := Register(tag, fn); err != nil {
        panic(err)
    }
}
