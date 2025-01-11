package apperrors

import "errors"

var (
	ErrAlreadyExists   = errors.New("value already exists")
	ErrNotFound        = errors.New("not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrDBQuery         = errors.New("database query error")
	ErrHashPassword    = errors.New("error hashing password")
	ErrJWTGeneration   = errors.New("error generating JWT")
)
