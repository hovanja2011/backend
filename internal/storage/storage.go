package storage

import "errors"

var (
	ErrUserNotFound = errors.New("url not found")
	ErrUserExists   = errors.New("url exists")
)
