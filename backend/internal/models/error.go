package models

import "errors"

var (
	ErrSessionEmpty     = errors.New("user session not found")
	ErrClientIPNotFound = errors.New("client ip not found")
	ErrToken            = errors.New("tokens do not match")

	ErrPassword   = errors.New("passwords do not match")
	ErrUsersEmpty = errors.New("user list is empty")
	ErrUserExist  = errors.New("user already exists")

	ErrOperationNotFound = errors.New("operation not found")
)
