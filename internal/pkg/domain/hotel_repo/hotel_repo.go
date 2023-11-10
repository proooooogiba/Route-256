package hotel_repo

import (
	"context"
	"errors"
	"github.com/opentracing/opentracing-go"
	"homework-3/internal/pkg/domain"
	"homework-3/internal/pkg/models"
	"homework-3/internal/pkg/repository"
)

type Hotel struct {
	db repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(db repository.DatabaseRepo) *Hotel {
	return &Hotel{
		db: db,
	}
}

func (m *Hotel) GetRoomWithAllReservations(ctx context.Context, id int64) (*models.Room, []*models.Reservation, error) {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"domain: get-room-with-all-reservations",
	)
	defer span.Finish()

	room, err := m.db.GetRoomByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, nil, domain.ErrRoomNotFound
		}
		return nil, nil, domain.ErrInternalServer
	}

	reservations, err := m.db.GetReservationsByRoomID(ctx, id)
	if err != nil {
		return nil, nil, domain.ErrInternalServer
	}
	return room, reservations, nil
}

func (m *Hotel) CreateRoom(ctx context.Context, room models.Room) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"domain: create-room",
	)
	defer span.Finish()

	_, err := m.db.GetRoomByName(ctx, room.Name)

	if err == nil {
		return 0, domain.ErrRoomAlreadyExists
	}

	if !errors.Is(err, repository.ErrObjectNotFound) {
		return 0, domain.ErrInternalServer
	}
	span.Finish()

	roomID, err := m.db.InsertRoom(ctx, &room)
	if err != nil {
		return 0, domain.ErrInternalServer
	}

	return roomID, nil
}

func (m *Hotel) UpdateRoom(ctx context.Context, room models.Room) error {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"domain: update-room",
	)
	defer span.Finish()

	_, err := m.db.GetRoomByID(ctx, room.ID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			//return http.StatusNotFound
			return domain.ErrRoomNotFound
		}
		return domain.ErrInternalServer
	}

	err = m.db.UpdateRoom(ctx, &room)
	if err != nil {
		return domain.ErrInternalServer
		//return http.StatusInternalServerError
	}

	return nil
}

func (m *Hotel) DeleteRoomWithAllReservations(ctx context.Context, id int64) error {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"domain: delete-room-with-all-reservations",
	)
	defer span.Finish()

	err := m.db.DeleteRoomByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return domain.ErrRoomNotFound
		}
		return domain.ErrInternalServer
	}

	err = m.db.DeleteReservationsByRoomID(ctx, id)
	if err != nil {
		return domain.ErrInternalServer
	}

	return nil
}

func (m *Hotel) GetReservation(ctx context.Context, key int64) (*models.Reservation, error) {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"domain: get-reservation",
	)
	defer span.Finish()

	res, err := m.db.GetReservationByID(ctx, key)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, domain.ErrReservationNotFound
		}
		return nil, domain.ErrInternalServer
	}
	return res, nil
}

func (m *Hotel) CreateReservation(ctx context.Context, res models.Reservation) (int64, error) {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"domain: create-reservation",
	)
	defer span.Finish()

	_, err := m.db.GetRoomByID(ctx, res.RoomID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return 0, domain.ErrRoomNotFound
		}
		return 0, domain.ErrInternalServer
	}

	resID, err := m.db.InsertReservation(ctx, &res)
	if err != nil {
		return 0, domain.ErrInternalServer
	}
	return resID, nil
}

func (m *Hotel) UpdateReservation(ctx context.Context, res models.Reservation) error {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"domain: update-reservation",
	)
	defer span.Finish()

	_, err := m.db.GetReservationByID(ctx, res.ID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return domain.ErrReservationNotFound
		}
		return domain.ErrInternalServer
	}

	err = m.db.UpdateReservation(ctx, &res)
	if err != nil {
		return domain.ErrInternalServer
	}

	return nil
}

func (m *Hotel) DeleteReservation(ctx context.Context, id int64) error {
	span, ctx := opentracing.StartSpanFromContext(
		ctx,
		"domain: delete-reservation",
	)
	defer span.Finish()

	err := m.db.DeleteReservationByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return domain.ErrReservationNotFound
		}
		return domain.ErrInternalServer
	}
	return nil
}
