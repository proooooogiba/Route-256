package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"homework-3/internal/handlers"
	"net/http"
)

func routes(repo *handlers.Hotel) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/room", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			body, status := handlers.GetBodyFromRequest(r)
			if status != http.StatusOK {
				w.WriteHeader(status)
				return
			}
			status = repo.CreateRoom(body)
			if status != http.StatusOK {
				w.WriteHeader(status)
				return
			}
			w.WriteHeader(status)
		case http.MethodPut:
			body, status := handlers.GetBodyFromRequest(r)
			if status != http.StatusOK {
				w.WriteHeader(status)
				return
			}
			status = repo.UpdateRoom(body)
			if status != http.StatusOK {
				w.WriteHeader(status)
				return
			}
			w.WriteHeader(status)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Method not allowed")
		}
	})

	router.HandleFunc(fmt.Sprintf("/room/{id:[0-9]+}"), func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			id, status := handlers.ParseGetID(r)
			if status != http.StatusOK {
				w.WriteHeader(status)
				return
			}
			jsonRoom, jsonReservations, status := repo.GetRoomWithAllReservations(id)
			w.WriteHeader(status)
			w.Write(jsonRoom)
			w.Write(jsonReservations)
		case http.MethodDelete:
			id, status := handlers.ParseGetID(r)
			if status != http.StatusOK {
				w.WriteHeader(status)
				return
			}
			status = repo.DeleteRoomWithAllReservations(id)
			w.WriteHeader(status)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Method not allowed")
		}
	})

	router.HandleFunc(fmt.Sprintf("/reservation"), func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			body, status := handlers.GetBodyFromRequest(r)
			if status != http.StatusOK {
				w.WriteHeader(status)
				return
			}
			status = repo.CreateReservation(body)
			w.WriteHeader(status)
		case http.MethodPut:
			body, status := handlers.GetBodyFromRequest(r)
			if status != http.StatusOK {
				w.WriteHeader(status)
				return
			}
			status = repo.UpdateReservation(body)
			w.WriteHeader(status)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Method not allowed")
		}
	})

	router.HandleFunc(fmt.Sprintf("/reservation/{id:[0-9]}"), func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			key, status := handlers.ParseGetID(r)
			if status != http.StatusOK {
				w.WriteHeader(status)
				return
			}
			data, status := repo.GetReservation(key)
			w.WriteHeader(status)
			w.Write(data)
		case http.MethodDelete:
			key, status := handlers.ParseGetID(r)
			if status != http.StatusOK {
				w.WriteHeader(status)
				return
			}
			status = repo.DeleteReservation(key)
			w.WriteHeader(status)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Method not allowed")
		}
	})

	return router
}
