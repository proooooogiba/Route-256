package handlers

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"homework-3/internal/pkg/db"
	"homework-3/internal/pkg/models"
	"homework-3/internal/pkg/repository"
	"homework-3/internal/pkg/repository/dbrepo"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Repository struct {
	DB repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(db *db.Database) *Repository {
	return &Repository{
		DB: dbrepo.NewPostgresRepo(db),
	}
}

type createRoomRequest struct {
	Name string  `json:"name"`
	Cost float64 `json:"cost"`
}

type updateRoomRequest struct {
	ID   int64   `json:"id"`
	Name string  `json:"name"`
	Cost float64 `json:"cost"`
}

type createReservationRequest struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	RoomID    int64  `json:"room_id"`
}

type updateReservationRequest struct {
	ID        int64  `json:"id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	RoomID    int64  `json:"room_id"`
}

func (m *Repository) GetRoomWithAllReservations(w http.ResponseWriter, r *http.Request) {
	IDStr, ok := mux.Vars(r)["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(IDStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	room, err := m.DB.GetRoomByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	reservations, err := m.DB.GetReservationsByRoomID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	roomJson, err := json.Marshal(room)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(roomJson)
	if len(reservations) == 0 {
		w.WriteHeader(http.StatusOK)
		return
	}

	for _, res := range reservations {
		w.Write([]byte("\n"))
		resJson, err := json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(resJson)
	}
	w.WriteHeader(http.StatusOK)
}

func (m *Repository) CreateRoom(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var unm createRoomRequest
	if err = json.Unmarshal(body, &unm); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	room := models.Room{
		Name:      unm.Name,
		Cost:      unm.Cost,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = m.DB.GetRoomByName(room.Name)

	if err == nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	if !errors.Is(err, repository.ErrObjectNotFound) {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = m.DB.InsertRoom(room)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (m *Repository) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var unm updateRoomRequest
	if err = json.Unmarshal(body, &unm); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	room := models.Room{
		ID:        unm.ID,
		Name:      unm.Name,
		Cost:      unm.Cost,
		UpdatedAt: time.Now(),
	}

	_, err = m.DB.GetRoomByID(room.ID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = m.DB.UpdateRoom(room)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (m *Repository) DeleteRoomWithAllReservations(w http.ResponseWriter, r *http.Request) {
	IDStr, ok := mux.Vars(r)["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(IDStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = m.DB.DeleteRoomByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = m.DB.DeleteReservationsByRoomID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (m *Repository) GetReservation(w http.ResponseWriter, r *http.Request) {
	IDStr, ok := mux.Vars(r)["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(IDStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := m.DB.GetReservationByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resJson, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(resJson)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (m *Repository) CreateReservation(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var unm createReservationRequest
	if err = json.Unmarshal(body, &unm); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, unm.StartDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	endDate, err := time.Parse(layout, unm.EndDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    unm.RoomID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = m.DB.InsertReservation(res)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (m *Repository) DeleteReservation(w http.ResponseWriter, r *http.Request) {
	IDStr, ok := mux.Vars(r)["id"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(IDStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = m.DB.DeleteReservationByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (m *Repository) UpdateReservation(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var unm updateReservationRequest
	if err = json.Unmarshal(body, &unm); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, unm.StartDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	endDate, err := time.Parse(layout, unm.EndDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res := models.Reservation{
		ID:        unm.ID,
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    unm.RoomID,
		UpdatedAt: time.Now(),
	}

	_, err = m.DB.GetReservationByID(res.ID)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = m.DB.UpdateReservation(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
