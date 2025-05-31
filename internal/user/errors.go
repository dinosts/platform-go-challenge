package user

import "errors"

var (
	ErrLoginFailed           = errors.New("Failed to login")
	ErrTokenGenerationFailed = errors.New("Could not generate jwtoken for user.")
	ErrUserNotFound          = errors.New("User Not Found")
)
