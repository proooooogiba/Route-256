//go:generate mockgen -source ./service.go -destination=./mocks/service.go -package=mock_service

package producer

import (
	"homework-3/internal/handlers"
	"homework-3/internal/pkg/models"
	"homework-3/internal/pkg/parser"
	"homework-3/internal/pkg/sender"
)

type Service struct {
	repo   handlers.Repository
	sender sender.Sender
	parser parser.Parser
}

func NewService(repo handlers.Repository, sender sender.Sender, parser parser.Parser) *Service {
	return &Service{
		repo:   repo,
		sender: sender,
		parser: parser,
	}
}

func (s Service) GetRoomWithAllReservations(roomID int64, sync bool) (*models.Room, []*models.Reservation, error) {
	room, reservations, err := s.repo.GetRoomWithAllReservations(roomID)
	if err != nil {
		return nil, nil, err
	}
	err = s.sender.Send("GET", []byte(""), sync)
	if err != nil {
		return nil, nil, err
	}

	return room, reservations, nil
}

func (s Service) CreateRoom(body []byte, sync bool) error {
	room, err := s.parser.UnmarshalCreateRoomRequest(body)
	if err != nil {
		return err
	}

	err = s.repo.CreateRoom(room)
	if err != nil {
		return err
	}

	err = s.sender.Send("POST", body, sync)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) UpdateRoom(body []byte, sync bool) error {
	room, err := s.parser.UnmarshalUpdateRoomRequest(body)
	if err != nil {
		return err
	}
	err = s.repo.UpdateRoom(room)
	if err != nil {
		return err
	}

	err = s.sender.Send("PUT", body, sync)
	if err != nil {
		return err
	}
	return nil
}

func (s Service) DeleteRoomWithAllReservations(roomID int64, sync bool) error {
	err := s.repo.DeleteRoomWithAllReservations(roomID)
	if err != nil {
		return err
	}
	err = s.sender.Send("DELETE", []byte(""), sync)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) GetReservation(resID int64, sync bool) (*models.Reservation, error) {
	reservation, err := s.repo.GetReservation(resID)
	if err != nil {
		return nil, err
	}
	err = s.sender.Send("GET", []byte(""), sync)
	if err != nil {
		return nil, err
	}

	return reservation, nil
}

func (s Service) DeleteReservation(resID int64, sync bool) error {
	err := s.repo.DeleteReservation(resID)
	if err != nil {
		return err
	}
	err = s.sender.Send("DELETE", []byte(""), sync)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) CreateReservation(body []byte, sync bool) error {
	reservation, err := s.parser.UnmarshalCreateReservationRequest(body)
	if err != nil {
		return err
	}

	err = s.repo.CreateReservation(reservation)
	if err != nil {
		return err
	}

	err = s.sender.Send("POST", body, sync)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) UpdateReservation(body []byte, sync bool) error {
	reservation, err := s.parser.UnmarshalUpdateReservationRequest(body)
	if err != nil {
		return err
	}

	err = s.repo.UpdateReservation(reservation)
	if err != nil {
		return err
	}

	err = s.sender.Send("PUT", body, sync)
	if err != nil {
		return err
	}

	return nil
}
