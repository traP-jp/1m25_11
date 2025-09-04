package repository

import "errors"

var (
	ErrNotFound = errors.New("entity not found")

	ErrAlreadyExists = errors.New("entity already exists")

	ErrForbidden = errors.New("operation not permitted")
)
