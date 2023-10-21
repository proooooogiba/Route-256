package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"homework-3/internal/handlers"
	"homework-3/internal/pkg/models"
	"io"
	"net/http"
	"strconv"
	"time"
)

func routes(hotel *handlers.Hotel) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/room", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			body, status := GetBodyFromRequest(r)
			if status != http.StatusOK {
				w.WriteHeader(status)
				return
			}
			room, status := UnmarshalCreateRoomRequest(body)
			if status != http.StatusOK {
				w.WriteHeader(status)
				return
			}

			err := hotel.CreateRoom(room)
			switch {
			case errors.Is(err, nil):
				w.WriteHeader(http.StatusOK)
			case errors.Is(err, handlers.ErrInternalServer):
				w.WriteHeader(http.StatusInternalServerError)
			case errors.Is(err, handlers.ErrRoomAlreadyExists):
				w.WriteHeader(http.StatusConflict)
			}
		case http.MethodPut:
			body, status := GetBodyFromRequest(r)
			if status != http.StatusOK {
				w.WriteHeader(status)
				return
			}

			room, status := UnmarshalUpdateRoomRequest(body)
			if status != http.StatusOK {
				w.WriteHeader(status)
				return
			}

			err := hotel.UpdateRoom(room)
			switch {
			case errors.Is(err, nil):
				w.WriteHeader(http.StatusOK)
			case errors.Is(err, handlers.ErrInternalServer):
				w.WriteHeader(http.StatusInternalServerError)
			case errors.Is(err, handlers.ErrRoomNotFound):
				w.WriteHeader(http.StatusNotFound)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Method not allowed")
		}
	})

	router.HandleFunc(fmt.Sprintf("/room/{id:[0-9]+}"), func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			id, status := ParseGetID(r)
			if status != http.StatusOK {
				w.WriteHeader(status)
				return
			}

			room, reservations, err := hotel.GetRoomWithAllReservations(id)
			if err != nil {
				switch {
				case errors.Is(err, nil):
					w.WriteHeader(http.StatusOK)
				case errors.Is(err, handlers.ErrInternalServer):
					w.WriteHeader(http.StatusInternalServerError)
				case errors.Is(err, handlers.ErrRoomNotFound):
					w.WriteHeader(http.StatusNotFound)
				}
				return
			}

			jsonRoom, err := json.Marshal(room)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			allResJson := bytes.NewBuffer([]byte("\n"))

			for _, res := range reservations {
				resJson, err := json.Marshal(res)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				allResJson.Write(resJson)
				allResJson.WriteByte('\n')
			}

			w.WriteHeader(http.StatusOK)
			w.Write(jsonRoom)
			w.Write(allResJson.Bytes())
		case http.MethodDelete:
			id, status := ParseGetID(r)
			if status != http.StatusOK {
				w.WriteHeader(status)
				return
			}
			err := hotel.DeleteRoomWithAllReservations(id)
			switch {
			case errors.Is(err, nil):
				w.WriteHeader(http.StatusOK)
			case errors.Is(err, handlers.ErrInternalServer):
				w.WriteHeader(http.StatusInternalServerError)
			case errors.Is(err, handlers.ErrReservationNotFound):
				w.WriteHeader(http.StatusNotFound)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Method not allowed")
		}
	})

	router.HandleFunc(fmt.Sprintf("/reservation"), func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			body, status := GetBodyFromRequest(r)
			if status != http.StatusOK {
				w.WriteHeader(status)
				return
			}

			res, status := UnmarshalCreateReservationRequest(body)
			if status != http.StatusOK {
				w.WriteHeader(status)
				return
			}

			err := hotel.CreateReservation(res)
			switch {
			case errors.Is(err, nil):
				w.WriteHeader(http.StatusOK)
			case errors.Is(err, handlers.ErrInternalServer):
				w.WriteHeader(http.StatusInternalServerError)
			case errors.Is(err, handlers.ErrRoomNotFound):
				w.WriteHeader(http.StatusNotFound)
			}
		case http.MethodPut:
			body, status := GetBodyFromRequest(r)
			if status != http.StatusOK {
				w.WriteHeader(status)
				return
			}

			res, status := UnmarshalUpdateReservationRequest(body)
			if status != http.StatusOK {
				w.WriteHeader(status)
				return
			}

			err := hotel.UpdateReservation(res)
			switch {
			case errors.Is(err, nil):
				w.WriteHeader(http.StatusOK)
			case errors.Is(err, handlers.ErrInternalServer):
				w.WriteHeader(http.StatusInternalServerError)
			case errors.Is(err, handlers.ErrReservationNotFound):
				w.WriteHeader(http.StatusNotFound)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Method nhttp.StatusOKot allowed")
		}
	})

	router.HandleFunc(fmt.Sprintf("/reservation/{id:[0-9]}"), func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			key, status := ParseGetID(r)
			if status != http.StatusOK {
				w.WriteHeader(status)
				return
			}
			res, err := hotel.GetReservation(key)
			if err != nil {
				switch {
				case errors.Is(err, handlers.ErrReservationNotFound):
					w.WriteHeader(http.StatusNotFound)
				case errors.Is(err, handlers.ErrInternalServer):
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			}

			resJson, err := json.Marshal(res)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write(resJson)
		case http.MethodDelete:
			key, status := ParseGetID(r)
			if status != http.StatusOK {
				w.WriteHeader(status)
				return
			}

			err := hotel.DeleteReservation(key)
			switch {
			case errors.Is(err, nil):
				w.WriteHeader(http.StatusOK)
			case errors.Is(err, handlers.ErrInternalServer):
				w.WriteHeader(http.StatusInternalServerError)
			case errors.Is(err, handlers.ErrReservationNotFound):
				w.WriteHeader(http.StatusNotFound)
			}

			w.WriteHeader(status)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Method not allowed")
		}
	})

	return router
}

func GetBodyFromRequest(r *http.Request) ([]byte, int) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	return body, http.StatusOK
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

func UnmarshalUpdateReservationRequest(body []byte) (models.Reservation, int) {
	var unm updateReservationRequest
	if err := json.Unmarshal(body, &unm); err != nil {
		return models.Reservation{}, http.StatusInternalServerError
	}

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, unm.StartDate)
	if err != nil {
		return models.Reservation{}, http.StatusInternalServerError
	}
	endDate, err := time.Parse(layout, unm.EndDate)
	if err != nil {
		return models.Reservation{}, http.StatusInternalServerError
	}

	res := models.Reservation{
		ID:        unm.ID,
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    unm.RoomID,
		UpdatedAt: time.Now(),
	}
	return res, http.StatusOK
}

func UnmarshalCreateReservationRequest(body []byte) (models.Reservation, int) {
	var unm CreateReservationRequest
	if err := json.Unmarshal(body, &unm); err != nil {
		return models.Reservation{}, http.StatusInternalServerError
	}

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, unm.StartDate)
	if err != nil {
		return models.Reservation{}, http.StatusInternalServerError
	}
	endDate, err := time.Parse(layout, unm.EndDate)
	if err != nil {
		return models.Reservation{}, http.StatusInternalServerError
	}

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    unm.RoomID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return res, http.StatusOK
}

func UnmarshalCreateRoomRequest(body []byte) (models.Room, int) {
	var unm CreateRoomRequest
	if err := json.Unmarshal(body, &unm); err != nil {
		return models.Room{}, http.StatusInternalServerError
	}

	room := models.Room{
		Name:      unm.Name,
		Cost:      unm.Cost,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return room, http.StatusOK
}

func UnmarshalUpdateRoomRequest(body []byte) (models.Room, int) {
	var unm UpdateRoomRequest
	if err := json.Unmarshal(body, &unm); err != nil {
		return models.Room{}, http.StatusInternalServerError
	}

	room := models.Room{
		ID:        unm.ID,
		Name:      unm.Name,
		Cost:      unm.Cost,
		UpdatedAt: time.Now(),
	}
	return room, http.StatusOK
}
