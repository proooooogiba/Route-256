package repository

import "errors"

var (
	ErrObjectNotFound   = errors.New("object not found")
	ErrObjectNotDelete  = errors.New("object not deleted")
	ErrObjectNotUpdated = errors.New("object not updated")
)
