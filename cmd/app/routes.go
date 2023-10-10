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
			repo.CreateRoom(w, r)
		case http.MethodPut:
			repo.UpdateRoom(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Method not allowed")
		}
	})

	router.HandleFunc(fmt.Sprintf("/room/{id:[0-9]+}"), func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			repo.GetRoomWithAllReservations(w, r)
		case http.MethodDelete:
			repo.DeleteRoomWithAllReservations(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Method not allowed")
		}
	})

	router.HandleFunc(fmt.Sprintf("/reservation"), func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			repo.CreateReservation(w, r)
		case http.MethodPut:
			repo.UpdateReservation(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Method not allowed")
		}
	})

	router.HandleFunc(fmt.Sprintf("/reservation/{id:[0-9]}"), func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			repo.GetReservation(w, r)
		case http.MethodDelete:
			repo.DeleteReservation(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Method not allowed")
		}
	})

	return router
}
