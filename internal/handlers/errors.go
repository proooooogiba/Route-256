package handlers

import "errors"

var (
	ErrObjectNotFound      = errors.New("object not found")
	ErrObjectNotDelete     = errors.New("object not deleted")
	ErrObjectNotUpdated    = errors.New("object not updated")
	ErrInternalServer      = errors.New("internal server error")
	ErrRoomAlreadyExists   = errors.New("room with same name already exists")
	ErrRoomNotFound        = errors.New("room with following id doesn't exists")
	ErrReservationNotFound = errors.New("reservation with following id doesn't exists")
	ErrSendSyncMessage     = errors.New("send sync message error")
	ErrSendASyncMessage    = errors.New("send async message error")
)
