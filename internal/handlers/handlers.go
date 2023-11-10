//go:generate mockgen -source ./service.go -destination=./mocks/service.go -package=mock_service

package handlers

import (
	"homework-3/internal/pkg/domain"
	"homework-3/internal/pkg/models"
	"homework-3/internal/pkg/parser"
	"homework-3/internal/pkg/sender"
)

type Service struct {
	repo   domain.Repository
	sender sender.Sender
	parser parser.Parser
}

func NewService(repo domain.Repository, sender sender.Sender, parser parser.Parser) *Service {
	return &Service{
		repo:   repo,
		sender: sender,
		parser: parser,
	}
}

func (s *Service) GetRoomWithAllReservations(method string, roomID int64, sync bool) (*models.Room, []*models.Reservation, error) {
	err := s.sender.Send(method, []byte(""), sync)
	if err != nil {
		return nil, nil, err
	}
	room, reservations, err := s.repo.GetRoomWithAllReservations(roomID)
	if err != nil {
		return nil, nil, err
	}

	return room, reservations, nil
}

func (s *Service) CreateRoom(method string, body []byte, sync bool) (int64, error) {
	room, err := s.parser.UnmarshalCreateRoomRequest(body)
	if err != nil {
		return 0, err
	}

	err = s.sender.Send(method, body, sync)
	if err != nil {
		return 0, err
	}

	roomID, err := s.repo.CreateRoom(room)
	if err != nil {
		return 0, err
	}

	return roomID, nil
}

func (s *Service) UpdateRoom(method string, body []byte, sync bool) error {
	room, err := s.parser.UnmarshalUpdateRoomRequest(body)
	if err != nil {
		return err
	}

	err = s.sender.Send(method, body, sync)
	if err != nil {
		return err
	}

	err = s.repo.UpdateRoom(room)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteRoomWithAllReservations(method string, roomID int64, sync bool) error {
	err := s.sender.Send(method, []byte(""), sync)
	if err != nil {
		return err
	}

	err = s.repo.DeleteRoomWithAllReservations(roomID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetReservation(method string, resID int64, sync bool) (*models.Reservation, error) {
	err := s.sender.Send(method, []byte(""), sync)
	if err != nil {
		return nil, err
	}

	reservation, err := s.repo.GetReservation(resID)
	if err != nil {
		return nil, err
	}

	return reservation, nil
}

func (s *Service) DeleteReservation(method string, resID int64, sync bool) error {
	err := s.sender.Send(method, []byte(""), sync)
	if err != nil {
		return err
	}

	err = s.repo.DeleteReservation(resID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) CreateReservation(method string, body []byte, sync bool) (int64, error) {
	reservation, err := s.parser.UnmarshalCreateReservationRequest(body)
	if err != nil {
		return 0, err
	}

	err = s.sender.Send(method, body, sync)
	if err != nil {
		return 0, err
	}

	resID, err := s.repo.CreateReservation(reservation)
	if err != nil {
		return 0, err
	}

	return resID, nil
}

func (s *Service) UpdateReservation(method string, body []byte, sync bool) error {
	reservation, err := s.parser.UnmarshalUpdateReservationRequest(body)
	if err != nil {
		return err
	}

	err = s.sender.Send(method, body, sync)
	if err != nil {
		return err
	}

	err = s.repo.UpdateReservation(reservation)
	if err != nil {
		return err
	}

	return nil
}
