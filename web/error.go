package web

// HTTPError as a response, with data
type HTTPError struct {
	Message string      `json:"message" example:"invalid format"`
	Error   string      `json:"error,omitempty" example:"<internal error message>"`
	Data    interface{} `json:"data,omitempty" swaggerignore:"true"`
}

const (
	// APIErrorInvalidRequestFormat const for api error with invalid request format
	APIErrorInvalidRequestFormat = "Invalid Request Format"
	// APIErrorInternalDatabase const for api error with problem with database
	APIErrorInternalDatabase = "Internal Database Error"
	// APIErrorNotFound const for api error with not found object
	APIErrorNotFound = "Not found"
)
