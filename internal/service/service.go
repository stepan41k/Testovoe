package service

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrSongExists         = errors.New("song already exists")
	ErrSongNotFound       = errors.New("song not found")
)