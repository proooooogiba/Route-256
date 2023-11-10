//go:generate mockgen -source ./repository.go -destination=./mocks/repository.go -package=mock_repo
package domain

import (
	"context"
	"homework-3/internal/pkg/models"
)

type Repository interface {
	GetRoomWithAllReservations(ctx context.Context, id int64) (*models.Room, []*models.Reservation, error)
	CreateRoom(ctx context.Context, room models.Room) (int64, error)
	UpdateRoom(ctx context.Context, room models.Room) error
	DeleteRoomWithAllReservations(ctx context.Context, id int64) error
	GetReservation(ctx context.Context, key int64) (*models.Reservation, error)
	CreateReservation(ctx context.Context, res models.Reservation) (int64, error)
	DeleteReservation(ctx context.Context, id int64) error
	UpdateReservation(ctx context.Context, res models.Reservation) error
}
