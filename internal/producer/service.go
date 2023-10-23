package producer

import (
	"encoding/json"
	"homework-3/internal/handlers"
	"homework-3/internal/pkg/models"
	"time"
)

type Sender interface {
	sendAsyncMessage(message RequestMessage) error
	sendMessage(message RequestMessage) error
	sendMessages(messages []RequestMessage) error
}

type Service struct {
	repo   handlers.Repository
	sender Sender
}

func NewService(repo handlers.Repository, sender Sender) *Service {
	return &Service{
		repo:   repo,
		sender: sender,
	}
}

func (s Service) GetRoomWithAllReservations(roomID int64, sync bool) (*models.Room, []*models.Reservation, error) {
	room, reservations, err := s.repo.GetRoomWithAllReservations(roomID)
	if err != nil {
		return nil, nil, err
	}
	err = s.Send("GET", []byte(""), sync)
	if err != nil {
		return nil, nil, err
	}

	return room, reservations, nil
}

func (s Service) CreateRoom(body []byte, sync bool) error {
	room, err := UnmarshalCreateRoomRequest(body)
	if err != nil {
		return err
	}

	err = s.Send("POST", body, sync)
	if err != nil {
		return err
	}

	err = s.repo.CreateRoom(room)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) UpdateRoom(body []byte, sync bool) error {
	room, err := UnmarshalUpdateRoomRequest(body)
	if err != nil {
		return err
	}
	err = s.repo.UpdateRoom(room)
	if err != nil {
		return err
	}

	err = s.Send("PUT", body, sync)
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
	err = s.Send("DELETE", []byte(""), sync)
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
	err = s.Send("GET", []byte(""), sync)
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
	err = s.Send("DELETE", []byte(""), sync)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) CreateReservation(body []byte, sync bool) error {
	reservation, err := UnmarshalCreateReservationRequest(body)
	if err != nil {
		return err
	}

	err = s.repo.CreateReservation(reservation)
	if err != nil {
		return err
	}

	err = s.Send("POST", body, sync)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) UpdateReservation(body []byte, sync bool) error {
	reservation, err := UnmarshalUpdateReservationRequest(body)
	if err != nil {
		return err
	}

	err = s.repo.UpdateReservation(reservation)
	if err != nil {
		return err
	}

	err = s.Send("PUT", body, sync)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) Send(method string, body []byte, sync bool) error {
	reqMsg := RequestMessage{
		Time:   time.Now(),
		Method: method,
		Body:   string(body),
	}
	switch sync {
	case true:
		err := s.sender.sendMessage(reqMsg)
		if err != nil {
			return handlers.ErrSendSyncMessage
		}
	case false:
		err := s.sender.sendAsyncMessage(reqMsg)

		if err != nil {
			return handlers.ErrSendASyncMessage
		}
	}
	return nil
}

func UnmarshalCreateRoomRequest(body []byte) (models.Room, error) {
	var unm createRoomRequest
	if err := json.Unmarshal(body, &unm); err != nil {
		return models.Room{}, err
	}

	room := models.Room{
		Name:      unm.Name,
		Cost:      unm.Cost,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return room, nil
}

func UnmarshalUpdateReservationRequest(body []byte) (models.Reservation, error) {
	var unm updateReservationRequest
	if err := json.Unmarshal(body, &unm); err != nil {
		return models.Reservation{}, err
	}

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, unm.StartDate)
	if err != nil {
		return models.Reservation{}, err
	}
	endDate, err := time.Parse(layout, unm.EndDate)
	if err != nil {
		return models.Reservation{}, err
	}

	res := models.Reservation{
		ID:        unm.ID,
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    unm.RoomID,
		UpdatedAt: time.Now(),
	}
	return res, nil
}

func UnmarshalCreateReservationRequest(body []byte) (models.Reservation, error) {
	var unm createReservationRequest
	if err := json.Unmarshal(body, &unm); err != nil {
		return models.Reservation{}, err
	}

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, unm.StartDate)
	if err != nil {
		return models.Reservation{}, err
	}
	endDate, err := time.Parse(layout, unm.EndDate)
	if err != nil {
		return models.Reservation{}, err
	}

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    unm.RoomID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return res, nil
}

func UnmarshalUpdateRoomRequest(body []byte) (models.Room, error) {
	var unm updateRoomRequest
	if err := json.Unmarshal(body, &unm); err != nil {
		return models.Room{}, err
	}

	room := models.Room{
		ID:        unm.ID,
		Name:      unm.Name,
		Cost:      unm.Cost,
		UpdatedAt: time.Now(),
	}
	return room, nil
}
