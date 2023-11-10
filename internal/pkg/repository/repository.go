//go:generate mockgen -source ./repository.go -destination=./mocks/repository.go -package=mock_repository

package repository

import (
	"context"
	"homework-3/internal/pkg/models"
)

type DatabaseRepo interface {
	InsertReservation(ctx context.Context, reservation *models.Reservation) (int64, error)
	GetReservationByID(ctx context.Context, id int64) (*models.Reservation, error)
	DeleteReservationByID(ctx context.Context, id int64) error
	UpdateReservation(ctx context.Context, res *models.Reservation) error
	InsertRoom(ctx context.Context, room *models.Room) (int64, error)
	GetRoomByID(ctx context.Context, id int64) (*models.Room, error)
	DeleteRoomByID(ctx context.Context, id int64) error
	UpdateRoom(ctx context.Context, room *models.Room) error
	GetRoomByName(ctx context.Context, name string) (*models.Room, error)
	GetReservationsByRoomID(ctx context.Context, roomID int64) ([]*models.Reservation, error)
	DeleteReservationsByRoomID(ctx context.Context, roomId int64) error
}
