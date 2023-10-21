package handlers

import (
	"errors"
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

func (m *Hotel) GetRoomWithAllReservations(id int64) (*models.Room, []*models.Reservation, error) {
	room, err := m.db.GetRoomByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, nil, ErrRoomNotFound
		}
		return nil, nil, ErrInternalServer
	}

	reservations, err := m.db.GetReservationsByRoomID(id)
	if err != nil {
		return nil, nil, ErrInternalServer
	}
	return room, reservations, nil
}

func (m *Hotel) CreateRoom(room models.Room) error {
	_, err := m.db.GetRoomByName(room.Name)

	if err == nil {
		return ErrRoomAlreadyExists
	}

	if !errors.Is(err, repository.ErrObjectNotFound) {
		return ErrInternalServer
	}

	_, err = m.db.InsertRoom(&room)
	if err != nil {
		return ErrInternalServer
	}
	return nil
}

func (m *Hotel) UpdateRoom(room models.Room) error {
	_, err := m.db.GetRoomByID(room.ID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			//return http.StatusNotFound
			return ErrRoomNotFound
		}
		return ErrInternalServer
		//return http.StatusInternalServerError
	}

	err = m.db.UpdateRoom(&room)
	if err != nil {
		return ErrInternalServer
		//return http.StatusInternalServerError
	}

	return nil
}

func (m *Hotel) DeleteRoomWithAllReservations(id int64) error {
	err := m.db.DeleteRoomByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return ErrRoomNotFound
		}
		return ErrInternalServer
	}

	err = m.db.DeleteReservationsByRoomID(id)
	if err != nil {
		return ErrInternalServer
	}

	return nil
}

func (m *Hotel) GetReservation(key int64) (*models.Reservation, error) {
	res, err := m.db.GetReservationByID(key)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, ErrReservationNotFound
		}
		return nil, ErrInternalServer
	}
	return res, nil
}

func (m *Hotel) CreateReservation(res models.Reservation) error {
	_, err := m.db.GetRoomByID(res.RoomID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return ErrRoomNotFound
		}
		return ErrInternalServer
	}

	_, err = m.db.InsertReservation(&res)
	if err != nil {
		return ErrInternalServer
	}
	return nil
}

func (m *Hotel) DeleteReservation(id int64) error {
	err := m.db.DeleteReservationByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return ErrReservationNotFound
		}
		return ErrInternalServer
	}
	return nil
}

func (m *Hotel) UpdateReservation(res models.Reservation) error {
	_, err := m.db.GetReservationByID(res.ID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return ErrReservationNotFound
		}
		return ErrInternalServer
	}

	err = m.db.UpdateReservation(&res)
	if err != nil {
		return ErrInternalServer
	}

	return nil
}
