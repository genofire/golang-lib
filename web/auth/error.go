package auth

const (
	// APIErrorUserNotFound api error string if user not found
	APIErrorUserNotFound string = "user not found"
	// APIErrorIncorrectPassword api error string if password is incorrect
	APIErrorIncorrectPassword string = "incorrect password"
	// APIErrorNoSession api error string if no session exists
	APIErrorNoSession string = "no session"
	// APIErrorCreateSession api error string if session could not created
	APIErrorCreateSession string = "create session"

	// APIErrroCreatePassword api error string if password could not created
	APIErrroCreatePassword string = "error during create password"
)
