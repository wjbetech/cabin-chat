package store

import "errors"

var (
	ErrUserNotFound                = errors.New("user not found")
	ErrUserAlreadyExists           = errors.New("user already exists")
	ErrUserAlreadyReactedWithEmoji = errors.New("user has already reacted with this emoji")
)
