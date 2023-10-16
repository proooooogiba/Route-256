package fixtures

import (
	"homework-3/internal/pkg/models"
	"homework-3/tests/states"
	"time"
)

type ReservationBuilder struct {
	instance *models.Reservation
}

func Reservation() *ReservationBuilder {
	return &ReservationBuilder{instance: &models.Reservation{}}
}

func (b *ReservationBuilder) ID(v int64) *ReservationBuilder {
	b.instance.ID = v
	return b
}

func (b *ReservationBuilder) StartDate(v time.Time) *ReservationBuilder {
	b.instance.StartDate = v
	return b
}

func (b *ReservationBuilder) EndDate(v time.Time) *ReservationBuilder {
	b.instance.EndDate = v
	return b
}

func (b *ReservationBuilder) RoomID(v int64) *ReservationBuilder {
	b.instance.RoomID = v
	return b
}

func (b *ReservationBuilder) CreatedAt(v time.Time) *ReservationBuilder {
	b.instance.CreatedAt = v
	return b
}

func (b *ReservationBuilder) UpdatedAt(v time.Time) *ReservationBuilder {
	b.instance.CreatedAt = v
	return b
}

func (b *ReservationBuilder) P() *models.Reservation {
	return b.instance
}

func (b *ReservationBuilder) V() models.Reservation {
	return *b.instance
}

func (b *ReservationBuilder) Valid() *ReservationBuilder {
	return Reservation().
		ID(1).
		StartDate(states.Reservation1StartDate).
		EndDate(states.Reservation1EndDate).
		RoomID(1).
		CreatedAt(time.Time{}).
		UpdatedAt(time.Time{})
}

func (b *ReservationBuilder) Valid2() *ReservationBuilder {
	return Reservation().
		ID(2).
		StartDate(states.Reservation2StartDate).
		EndDate(states.Reservation2EndDate).
		RoomID(1).
		CreatedAt(time.Time{}).
		UpdatedAt(time.Time{})
}

func (b *RoomBuilder) ValidWithoutID() *ReservationBuilder {
	return Reservation().
		StartDate(states.Reservation1StartDate).
		EndDate(states.Reservation1EndDate).
		RoomID(1).
		CreatedAt(time.Time{}).
		UpdatedAt(time.Time{})
}
