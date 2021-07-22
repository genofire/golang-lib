package auth

import "errors"

var (
	// ErrAPIUserNotFound api error string if user not found
	ErrAPIUserNotFound = errors.New("user not found")
	// ErrAPIIncorrectPassword api error string if password is incorrect
	ErrAPIIncorrectPassword = errors.New("incorrect password")
	// ErrAPINoSession api error string if no session exists
	ErrAPINoSession = errors.New("no session")
	// ErrAPICreateSession api error string if session could not created
	ErrAPICreateSession = errors.New("create session")

	// ErrAPICreatePassword api error string if password could not created
	ErrAPICreatePassword = errors.New("error during create password")

	// ErrAPINoPermission api error string if an error happen on accesing this object
	ErrAPINoPermission = errors.New("error on access an object")
)
