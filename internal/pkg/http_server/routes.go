package http_server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"homework-3/internal/handlers"
	"homework-3/internal/pkg/domain"
	"io"
	"net/http"
	"strconv"
)

func Routes(hotel *handlers.Service) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/room", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			body, status := GetBodyFromRequest(r)
			if status != http.StatusOK {
				w.WriteHeader(status)
				return
			}
			_, err := hotel.CreateRoom(http.MethodPost, body, true)

			switch {
			case errors.Is(err, nil):
				w.WriteHeader(http.StatusOK)
			case errors.Is(err, domain.ErrInternalServer):
				w.WriteHeader(http.StatusInternalServerError)
			case errors.Is(err, domain.ErrRoomAlreadyExists):
				w.WriteHeader(http.StatusConflict)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
		case http.MethodPut:
			body, status := GetBodyFromRequest(r)
			if status != http.StatusOK {
				w.WriteHeader(status)
				return
			}

			err := hotel.UpdateRoom(http.MethodPut, body, true)
			switch {
			case errors.Is(err, nil):
				w.WriteHeader(http.StatusOK)
			case errors.Is(err, domain.ErrInternalServer):
				w.WriteHeader(http.StatusInternalServerError)
			case errors.Is(err, domain.ErrRoomNotFound):
				w.WriteHeader(http.StatusNotFound)
			default:
				w.WriteHeader(http.StatusInternalServerError)
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

			room, reservations, err := hotel.GetRoomWithAllReservations(http.MethodGet, id, true)
			if err != nil {
				switch {
				case errors.Is(err, nil):
					w.WriteHeader(http.StatusOK)
				case errors.Is(err, domain.ErrInternalServer):
					w.WriteHeader(http.StatusInternalServerError)
				case errors.Is(err, domain.ErrRoomNotFound):
					w.WriteHeader(http.StatusNotFound)
				default:
					w.WriteHeader(http.StatusInternalServerError)
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
			err := hotel.DeleteRoomWithAllReservations(http.MethodDelete, id, false)
			switch {
			case errors.Is(err, nil):
				w.WriteHeader(http.StatusOK)
			case errors.Is(err, domain.ErrInternalServer):
				w.WriteHeader(http.StatusInternalServerError)
			case errors.Is(err, domain.ErrReservationNotFound):
				w.WriteHeader(http.StatusNotFound)
			default:
				w.WriteHeader(http.StatusInternalServerError)
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

			_, err := hotel.CreateReservation(http.MethodPost, body, true)
			switch {
			case errors.Is(err, nil):
				w.WriteHeader(http.StatusOK)
			case errors.Is(err, domain.ErrInternalServer):
				w.WriteHeader(http.StatusInternalServerError)
			case errors.Is(err, domain.ErrRoomNotFound):
				w.WriteHeader(http.StatusNotFound)
			default:
				w.WriteHeader(http.StatusInternalServerError)
			}
		case http.MethodPut:
			body, status := GetBodyFromRequest(r)
			if status != http.StatusOK {
				w.WriteHeader(status)
				return
			}

			err := hotel.UpdateReservation(http.MethodPut, body, false)
			switch {
			case errors.Is(err, nil):
				w.WriteHeader(http.StatusOK)
			case errors.Is(err, domain.ErrInternalServer):
				w.WriteHeader(http.StatusInternalServerError)
			case errors.Is(err, domain.ErrReservationNotFound):
				w.WriteHeader(http.StatusNotFound)
			default:
				w.WriteHeader(http.StatusInternalServerError)
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
			res, err := hotel.GetReservation(http.MethodGet, key, true)
			if err != nil {
				switch {
				case errors.Is(err, domain.ErrReservationNotFound):
					w.WriteHeader(http.StatusNotFound)
				case errors.Is(err, domain.ErrInternalServer):
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

			err := hotel.DeleteReservation(http.MethodDelete, key, false)
			switch {
			case errors.Is(err, nil):
				w.WriteHeader(http.StatusOK)
			case errors.Is(err, domain.ErrInternalServer):
				w.WriteHeader(http.StatusInternalServerError)
			case errors.Is(err, domain.ErrReservationNotFound):
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
