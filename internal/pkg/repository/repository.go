//go:generate mockgen -source ./repository.go -destination=./mocks/repository.go -package=mock_repository

package repository

import (
	"homework-3/internal/pkg/models"
)

type DatabaseRepo interface {
	InsertReservation(reservation *models.Reservation) (int64, error)
	GetReservationByID(id int64) (*models.Reservation, error)
	DeleteReservationByID(id int64) error
	UpdateReservation(res *models.Reservation) error
	InsertRoom(room *models.Room) (int64, error)
	GetRoomByID(id int64) (*models.Room, error)
	DeleteRoomByID(id int64) error
	UpdateRoom(room *models.Room) error
	GetRoomByName(name string) (*models.Room, error)
	GetReservationsByRoomID(roomID int64) ([]*models.Reservation, error)
	DeleteReservationsByRoomID(roomId int64) error
}
