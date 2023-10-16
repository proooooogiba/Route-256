package states

import "time"

const (
	Room1ID        = 200001
	Reservation1ID = 500001
	Reservation2ID = 500002
	Reservation3ID = 500003
)

const Room1Name = "Lux"

var (
	Reservation1StartDate = time.Date(2023, 11, 8, 0, 0, 0, 0, time.UTC)
	Reservation1EndDate   = time.Date(2023, 11, 17, 0, 0, 0, 0, time.UTC)
	Reservation2StartDate = time.Date(2023, 10, 8, 0, 0, 0, 0, time.UTC)
	Reservation2EndDate   = time.Date(2023, 10, 17, 0, 0, 0, 0, time.UTC)
)
