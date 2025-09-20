package responses

import "encoding/json"

type (
	ErrorResponse struct {
		Error Error `json:"error"`
	}

	Error struct {
		Message string `json:"message"`
	}

	BoolResponse struct {
		Result bool `json:"result"`
	}
)

const (
	ErrorCodeDefault = 0
)

func newErrorResponse(message string) ErrorResponse {
	return ErrorResponse{
		Error: Error{
			Message: message,
		},
	}
}

func EncodeNotFoundResponse() ([]byte, error) {
	response := newErrorResponse("not found")
	return json.Marshal(response)
}

func EncodeServerErrorResponse() ([]byte, error) {
	response := newErrorResponse("server error")
	return json.Marshal(response)
}

func EncodeOkResponse() ([]byte, error) {
	return json.Marshal(BoolResponse{
		Result: true,
	})
}
