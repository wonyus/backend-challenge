package errors

import "errors"

var (
	// User errors
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidUserData    = errors.New("invalid user data")

	// Auth errors
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
	ErrUnauthorized = errors.New("unauthorized")

	// Validation errors
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrPasswordTooShort = errors.New("password too short")
	ErrRequiredField    = errors.New("required field missing")

	// Jwt errors
	ErrEmptyPassword      = errors.New("password cannot be empty")
	ErrInvalidTokenSecret = errors.New("invalid token secret")
)
