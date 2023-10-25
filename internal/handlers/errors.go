package handlers

import "errors"

var (
	ErrInternalServer      = errors.New("internal server error")
	ErrRoomAlreadyExists   = errors.New("room with same name already exists")
	ErrRoomNotFound        = errors.New("room with following id doesn't exists")
	ErrReservationNotFound = errors.New("reservation with following id doesn't exists")
)
