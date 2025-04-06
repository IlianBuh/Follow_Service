package storage

import "errors"

var (
	ErrFollowing   = errors.New("user is already following")
	ErrNoFollowing = errors.New("user has not followed")
)
