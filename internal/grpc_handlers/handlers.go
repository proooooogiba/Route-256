package grpc_handlers

import (
	"context"
	"homework-3/internal/pkg/domain"
	"homework-3/internal/pkg/models"
)

type Service struct {
	repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetRoomWithAllReservations(ctx context.Context, roomID int64) (*models.Room, []*models.Reservation, error) {
	room, reservations, err := s.repo.GetRoomWithAllReservations(roomID)
	if err != nil {
		return nil, nil, err
	}

	return room, reservations, err
}

func (s *Service) CreateRoom(ctx context.Context, room models.Room) (int64, error) {
	roomID, err := s.repo.CreateRoom(room)
	if err != nil {
		return 0, err
	}

	return roomID, err
}

func (s *Service) UpdateRoom(ctx context.Context, room models.Room) error {
	err := s.repo.UpdateRoom(room)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteRoomWithAllReservations(ctx context.Context, roomID int64) error {
	err := s.repo.DeleteRoomWithAllReservations(roomID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) CreateReservation(ctx context.Context, reservation models.Reservation) (int64, error) {
	resID, err := s.repo.CreateReservation(reservation)
	if err != nil {
		return 0, err
	}

	return resID, nil
}

func (s *Service) GetReservation(ctx context.Context, resID int64) (*models.Reservation, error) {
	reservation, err := s.repo.GetReservation(resID)
	if err != nil {
		return nil, err
	}

	return reservation, nil
}

func (s *Service) UpdateReservation(ctx context.Context, reservation models.Reservation) error {
	err := s.repo.UpdateReservation(reservation)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteReservation(ctx context.Context, resID int64) error {
	err := s.repo.DeleteReservation(resID)
	if err != nil {
		return err
	}

	return nil
}
