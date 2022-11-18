package models

import "errors"

var (
	ErrSessionEmpty     = errors.New("user session not found")
	ErrClientIPNotFound = errors.New("client ip not found")
	ErrToken            = errors.New("tokens do not match")
)
