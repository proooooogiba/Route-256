package main

type CreateRoomRequest struct {
	Name string  `json:"name"`
	Cost float64 `json:"cost"`
}

type UpdateRoomRequest struct {
	ID   int64   `json:"id"`
	Name string  `json:"name"`
	Cost float64 `json:"cost"`
}

type CreateReservationRequest struct {
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
