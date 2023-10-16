package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"homework-3/internal/pkg/models"
	"homework-3/internal/pkg/repository"
	"io"
	"net/http"
	"strconv"
	"time"
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

func (m *Hotel) GetRoomWithAllReservations(id int64) ([]byte, []byte, int) {
	room, err := m.db.GetRoomByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, nil, http.StatusNotFound
		}
		return nil, nil, http.StatusInternalServerError
	}

	reservations, err := m.db.GetReservationsByRoomID(id)
	if err != nil {
		return nil, nil, http.StatusInternalServerError
	}

	roomJson, err := json.Marshal(room)
	if err != nil {
		return nil, nil, http.StatusInternalServerError
	}

	if len(reservations) == 0 {
		return roomJson, nil, http.StatusOK
	}

	allResJson := bytes.NewBuffer([]byte("\n"))

	for _, res := range reservations {
		resJson, err := json.Marshal(res)
		if err != nil {
			return nil, nil, http.StatusInternalServerError
		}
		allResJson.Write(resJson)
		allResJson.WriteByte('\n')
	}

	return roomJson, allResJson.Bytes(), http.StatusOK
}

func GetBodyFromRequest(r *http.Request) ([]byte, int) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	return body, http.StatusOK
}

func (m *Hotel) CreateRoom(body []byte) int {
	var unm createRoomRequest
	if err := json.Unmarshal(body, &unm); err != nil {
		return http.StatusInternalServerError
	}

	room := models.Room{
		Name:      unm.Name,
		Cost:      unm.Cost,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := m.db.GetRoomByName(room.Name)

	if err == nil {
		return http.StatusConflict
	}

	if !errors.Is(err, repository.ErrObjectNotFound) {
		return http.StatusInternalServerError
	}

	_, err = m.db.InsertRoom(&room)
	if err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func (m *Hotel) UpdateRoom(body []byte) int {
	var unm updateRoomRequest
	if err := json.Unmarshal(body, &unm); err != nil {
		return http.StatusInternalServerError
	}

	room := models.Room{
		ID:        unm.ID,
		Name:      unm.Name,
		Cost:      unm.Cost,
		UpdatedAt: time.Now(),
	}

	_, err := m.db.GetRoomByID(room.ID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return http.StatusNotFound
		}
		return http.StatusInternalServerError
	}

	err = m.db.UpdateRoom(&room)
	if err != nil {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

func (m *Hotel) DeleteRoomWithAllReservations(id int64) int {
	err := m.db.DeleteRoomByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return http.StatusNotFound
		}
		return http.StatusInternalServerError
	}

	err = m.db.DeleteReservationsByRoomID(id)
	if err != nil {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

func (m *Hotel) GetReservation(key int64) ([]byte, int) {
	res, err := m.db.GetReservationByID(key)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return nil, http.StatusNotFound
		}
		return nil, http.StatusInternalServerError
	}

	resJson, err := json.Marshal(res)
	if err != nil {
		return nil, http.StatusInternalServerError
	}

	return resJson, http.StatusOK
}

func ParseGetID(req *http.Request) (int64, int) {
	key, ok := mux.Vars(req)["id"]
	if !ok {
		return 0, http.StatusBadRequest
	}
	keyInt, err := strconv.ParseInt(key, 10, 64)
	if err != nil {
		return 0, http.StatusBadRequest
	}
	return keyInt, http.StatusOK
}

func (m *Hotel) CreateReservation(body []byte) int {
	var unm createReservationRequest
	if err := json.Unmarshal(body, &unm); err != nil {
		return http.StatusInternalServerError
	}

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, unm.StartDate)
	if err != nil {
		return http.StatusInternalServerError
	}
	endDate, err := time.Parse(layout, unm.EndDate)
	if err != nil {
		return http.StatusInternalServerError
	}
	_, err = m.db.GetRoomByID(unm.RoomID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return http.StatusNotFound
		}
		return http.StatusInternalServerError
	}

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    unm.RoomID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = m.db.InsertReservation(&res)
	if err != nil {
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func (m *Hotel) DeleteReservation(id int64) int {
	err := m.db.DeleteReservationByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return http.StatusNotFound
		}
		return http.StatusInternalServerError
	}
	return http.StatusOK
}

func (m *Hotel) UpdateReservation(body []byte) int {
	var unm updateReservationRequest
	if err := json.Unmarshal(body, &unm); err != nil {
		return http.StatusInternalServerError
	}

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, unm.StartDate)
	if err != nil {
		return http.StatusInternalServerError
	}
	endDate, err := time.Parse(layout, unm.EndDate)
	if err != nil {
		return http.StatusInternalServerError
	}

	res := models.Reservation{
		ID:        unm.ID,
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    unm.RoomID,
		UpdatedAt: time.Now(),
	}

	_, err = m.db.GetReservationByID(res.ID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return http.StatusNotFound
		}
		return http.StatusInternalServerError
	}

	err = m.db.UpdateReservation(&res)
	if err != nil {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}
