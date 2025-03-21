package storage

import "errors"

var (
	ErrSongExists = errors.New("song already exists")
	ErrSongNotFound = errors.New("song not found")
	ErrNoChanges = errors.New("no changes")
)