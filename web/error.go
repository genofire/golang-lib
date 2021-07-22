package web

import "errors"

// HTTPError is returned in HTTP error responses.
type HTTPError struct {
	Message string      `json:"message" example:"invalid format"`
	Error   string      `json:"error,omitempty" example:"<internal error message>"`
	Data    interface{} `json:"data,omitempty" swaggerignore:"true"`
}

// Error strings used for HTTPError.Message.
var (
	ErrAPIInvalidRequestFormat = errors.New("Invalid Request Format")
	ErrAPIInternalDatabase     = errors.New("Internal Database Error")
	ErrAPINotFound             = errors.New("Not found")
)
