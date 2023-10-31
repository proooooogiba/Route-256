package parser_request

import (
	"encoding/json"
	"homework-3/internal/pkg/models"
	"homework-3/internal/pkg/parser"
	"time"
)

type RequestParser struct{}

func NewRequestParser() *RequestParser {
	return &RequestParser{}
}

func (p *RequestParser) UnmarshalCreateRoomRequest(body []byte) (models.Room, error) {
	var unm createRoomRequest
	if err := json.Unmarshal(body, &unm); err != nil {
		return models.Room{}, parser.ErrUnmarshal
	}

	room := models.Room{
		Name:      unm.Name,
		Cost:      unm.Cost,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return room, nil
}

func (p *RequestParser) UnmarshalUpdateReservationRequest(body []byte) (models.Reservation, error) {
	var unm updateReservationRequest
	if err := json.Unmarshal(body, &unm); err != nil {
		return models.Reservation{}, parser.ErrUnmarshal
	}

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, unm.StartDate)
	if err != nil {
		return models.Reservation{}, parser.ErrParseDate
	}
	endDate, err := time.Parse(layout, unm.EndDate)
	if err != nil {
		return models.Reservation{}, parser.ErrParseDate
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

func (p *RequestParser) UnmarshalCreateReservationRequest(body []byte) (models.Reservation, error) {
	var unm createReservationRequest
	if err := json.Unmarshal(body, &unm); err != nil {
		return models.Reservation{}, parser.ErrUnmarshal
	}

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, unm.StartDate)
	if err != nil {
		return models.Reservation{}, parser.ErrParseDate
	}
	endDate, err := time.Parse(layout, unm.EndDate)
	if err != nil {
		return models.Reservation{}, parser.ErrParseDate
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

func (p *RequestParser) UnmarshalUpdateRoomRequest(body []byte) (models.Room, error) {
	var unm updateRoomRequest
	if err := json.Unmarshal(body, &unm); err != nil {
		return models.Room{}, parser.ErrUnmarshal
	}

	room := models.Room{
		ID:        unm.ID,
		Name:      unm.Name,
		Cost:      unm.Cost,
		UpdatedAt: time.Now(),
	}
	return room, nil
}
