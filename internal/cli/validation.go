package cli

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Gets a custom error message for a validated field.
func msgForTag(fe validator.FieldError) string {
	switch fe.Field() {
	case "Threads":
		return "'Threads' option must be a number between 1 and 10, inclusive"
	case "URLs":
		return "'URLs' option must be a list of valid URLs"
	}
	return fe.Error()
}

// Handles errors produced by go-playground/validator.
func handleValidationErr(err error) error {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		//Collect all validation errors
		var errs []string
		for _, fe := range validationErrors {
			errs = append(errs, msgForTag(fe))
		}

		//Produce a new error
		return fmt.Errorf("encountered %d validation error(s):\n  %s", len(errs), strings.Join(errs, "\n  "))
	} else {
		return err
	}
}

// Registers custom validation handlers for go-playground/validator.
func registerCustomHandlers(v *validator.Validate) {
	//Custom validation for []string where each element is a URL
	err := v.RegisterValidation("urlslice", func(fl validator.FieldLevel) bool {
		urls, ok := fl.Field().Interface().([]string)
		if !ok {
			return false
		}
		for _, u := range urls {
			err := v.Var(u, "required,url")
			if err != nil {
				return false
			}
		}
		return true
	})
	if err != nil {
		panic(err) //Shouldn't be thrown under normal conditions
	}
}
