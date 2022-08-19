package errors

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type BaseError struct {
	Type     string   `json:"type"`
	Messsage string   `json:"messsage"`
	Errors   []string `json:"errors"`
}

type InvalidJSONError struct {
	BaseError
}

func NewInvalidJSONError() *InvalidJSONError {

	return &InvalidJSONError{
		BaseError: BaseError{
			Type:     "invalid.json.error",
			Messsage: "Incorrect json format. Please check your json for correctness !",
		},
	}
}

type BodyValidatonError struct {
	BaseError
}

func NewValidationError(typeOfData string, err validator.ValidationErrors) *BodyValidatonError {
	var errorsArray []string

	for _, err := range err {

		errorsArray = append(errorsArray, fmt.Sprintf(`Field '%s' must be required`, strings.ToLower(err.Field())))

	}

	return &BodyValidatonError{
		BaseError: BaseError{
			Type:     "invalid.data",
			Messsage: fmt.Sprintf("Incorrect %s format. Please check your json for correctness !", typeOfData),
			Errors:   errorsArray,
		},
	}
}

type InvalidQueryParametersError struct {
	BaseError
}

func NewInvalidQueryParametersError(msgError error) *InvalidQueryParametersError {
	return &InvalidQueryParametersError{
		BaseError: BaseError{
			Type:     "invalid.query.parameters",
			Messsage: "The query contains incorrect filling of the query string !",
			Errors:   []string{msgError.Error()},
		},
	}
}
