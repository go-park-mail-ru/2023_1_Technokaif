package models


// AUTH ERRORS
type UserAlreadyExistsError struct{}
type NoSuchUserError struct{}

func (e *UserAlreadyExistsError) Error() string {
	return "user already exists"
}

func (e *NoSuchUserError) Error() string {
	return "no such user"
}