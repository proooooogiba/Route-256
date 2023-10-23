package handlers

import "homework-3/internal/pkg/models"

type Repository interface {
	GetRoomWithAllReservations(id int64) (*models.Room, []*models.Reservation, error)
	CreateRoom(room models.Room) error
	UpdateRoom(room models.Room) error
	DeleteRoomWithAllReservations(id int64) error
	GetReservation(key int64) (*models.Reservation, error)
	CreateReservation(res models.Reservation) error
	DeleteReservation(id int64) error
	UpdateReservation(res models.Reservation) error
}
