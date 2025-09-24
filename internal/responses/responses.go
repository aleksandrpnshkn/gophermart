package responses

import (
	"encoding/json"
)

type (
	ErrorResponse struct {
		Error Error `json:"error"`
	}

	Error struct {
		Message       string         `json:"message"`
		InvalidFields []InvalidField `json:"invalid_fields,omitempty"`
	}

	ValidationError struct {
		Message       string         `json:"message"`
		InvalidFields []InvalidField `json:"invalid_fields"`
	}

	InvalidField struct {
		Field   string `json:"field"`
		Message string `json:"message"`
	}

	BoolResponse struct {
		Result bool `json:"result"`
	}
)

func newErrorResponse(message string) ErrorResponse {
	return ErrorResponse{
		Error: Error{
			Message: message,
		},
	}
}

func EncodeErrorResponse(message string) ([]byte, error) {
	response := newErrorResponse(message)
	return json.Marshal(response)
}

func EncodeValidationErrorResponse(message string, errors map[string]string) ([]byte, error) {
	invalidFields := []InvalidField{}

	for field, message := range errors {
		invalidFields = append(invalidFields, InvalidField{
			Field:   field,
			Message: message,
		})
	}

	response := ErrorResponse{
		Error: Error{
			Message:       message,
			InvalidFields: invalidFields,
		},
	}
	return json.Marshal(response)
}

func EncodeNotFoundResponse() ([]byte, error) {
	return EncodeErrorResponse("not found")
}

func EncodeServerErrorResponse() ([]byte, error) {
	return EncodeErrorResponse("server error")
}

func EncodeOkResponse() ([]byte, error) {
	return json.Marshal(BoolResponse{
		Result: true,
	})
}
