package fixtures

import (
	"homework-3/internal/pkg/models"
	"homework-3/tests/states"
	"time"
)

type RoomBuilder struct {
	instance *models.Room
}

func Room() *RoomBuilder {
	return &RoomBuilder{instance: &models.Room{}}
}

func (b *RoomBuilder) ID(v int64) *RoomBuilder {
	b.instance.ID = v
	return b
}

func (b *RoomBuilder) Name(v string) *RoomBuilder {
	b.instance.Name = v
	return b
}

func (b *RoomBuilder) Cost(v float64) *RoomBuilder {
	b.instance.Cost = v
	return b
}

func (b *RoomBuilder) CreatedAt(v time.Time) *RoomBuilder {
	b.instance.CreatedAt = v
	return b
}

func (b *RoomBuilder) UpdatedAt(v time.Time) *RoomBuilder {
	b.instance.CreatedAt = v
	return b
}

func (b *RoomBuilder) P() *models.Room {
	return b.instance
}

func (b *RoomBuilder) V() models.Room {
	return *b.instance
}

func (b *RoomBuilder) Valid() *RoomBuilder {
	return Room().
		ID(1).
		Name(states.Room1Name).
		Cost(1000.0).
		CreatedAt(time.Time{}).
		UpdatedAt(time.Time{})
}

func (b *RoomBuilder) UpdatedValid() *RoomBuilder {
	return Room().
		ID(1).
		Name(states.Room1Name).
		Cost(1200.0).
		CreatedAt(time.Time{}).
		UpdatedAt(time.Time{})
}
