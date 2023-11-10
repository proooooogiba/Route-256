package grpc_handlers

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"homework-3/internal/pkg/domain"
	"homework-3/internal/pkg/logger"
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
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"service: get-room-with-all-reservations",
	)
	defer span.Finish()

	ctx = logger.SetMethod(ctx, "GetRoomWithAllReservations")

	room, reservations, err := s.repo.GetRoomWithAllReservations(ctx, roomID)
	if err != nil {
		logger.Errorf(ctx, "getting error while getting room with all reservations: %s", err)
		return nil, nil, err
	}

	return room, reservations, err
}

func (s *Service) CreateRoom(ctx context.Context, room models.Room) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"service: create-room",
	)
	defer span.Finish()

	ctx = logger.SetMethod(ctx, "CreateRoom")

	roomID, err := s.repo.CreateRoom(ctx, room)
	if err != nil {
		logger.Errorf(ctx, "getting error while creating room: %s", err)
		return 0, err
	}

	return roomID, err
}

func (s *Service) UpdateRoom(ctx context.Context, room models.Room) error {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"service: update-room",
	)
	defer span.Finish()

	ctx = logger.SetMethod(ctx, "UpdateRoom")

	err := s.repo.UpdateRoom(ctx, room)
	if err != nil {
		logger.Errorf(ctx, "getting error while updating room: %s", err)
		return err
	}

	return nil
}

func (s *Service) DeleteRoomWithAllReservations(ctx context.Context, roomID int64) error {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"service: delete-room-with-all-reservations",
	)
	defer span.Finish()

	ctx = logger.SetMethod(ctx, "DeleteRoomWithAllReservations")

	err := s.repo.DeleteRoomWithAllReservations(ctx, roomID)
	if err != nil {
		logger.Errorf(ctx, "getting error while deleting room with all reservations: %s", err)
		return err
	}

	return nil
}

func (s *Service) CreateReservation(ctx context.Context, reservation models.Reservation) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"service: create-reservation",
	)
	defer span.Finish()

	ctx = logger.SetMethod(ctx, "CreateReservation")

	resID, err := s.repo.CreateReservation(ctx, reservation)
	if err != nil {
		logger.Errorf(ctx, "getting error while creating reservation: %s", err)
		return 0, err
	}

	return resID, nil
}

func (s *Service) GetReservation(ctx context.Context, resID int64) (*models.Reservation, error) {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"service: get-reservation",
	)
	defer span.Finish()

	ctx = logger.SetMethod(ctx, "GetReservation")

	reservation, err := s.repo.GetReservation(ctx, resID)
	if err != nil {
		logger.Errorf(ctx, "getting error while getting reservation: %s", err)
		return nil, err
	}

	return reservation, nil
}

func (s *Service) UpdateReservation(ctx context.Context, reservation models.Reservation) error {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"service: update-reservation",
	)
	defer span.Finish()

	ctx = logger.SetMethod(ctx, "UpdateReservation")

	err := s.repo.UpdateReservation(ctx, reservation)
	if err != nil {
		logger.Errorf(ctx, "getting error while updating reservation: %s", err)
		return err
	}
	return nil
}

func (s *Service) DeleteReservation(ctx context.Context, resID int64) error {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"service: delete-reservation",
	)
	defer span.Finish()

	ctx = logger.SetMethod(ctx, "DeleteReservation")

	err := s.repo.DeleteReservation(ctx, resID)
	if err != nil {
		logger.Errorf(ctx, "getting error while deleting reservation: %s", err)
		return err
	}

	return nil
}
