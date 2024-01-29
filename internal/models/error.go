package models

import "errors"

// repository error
var (
	ErrUserNotFound  = errors.New("user not found")
	ErrGameNotFound  = errors.New("game not found")
	ErrGameDuplicate = errors.New("game duplicate")
)

// game error
var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrGameNotAvailable    = errors.New("game not available")
	ErrBadCommand          = errors.New("bad command")
)
