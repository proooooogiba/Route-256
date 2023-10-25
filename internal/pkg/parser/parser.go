//go:generate mockgen -source ./parser.go -destination=./mocks/parser.go -package=mock_parser

package parser

import (
	"homework-3/internal/pkg/models"
)

type Parser interface {
	UnmarshalCreateRoomRequest(body []byte) (models.Room, error)
	UnmarshalUpdateReservationRequest(body []byte) (models.Reservation, error)
	UnmarshalCreateReservationRequest(body []byte) (models.Reservation, error)
	UnmarshalUpdateRoomRequest(body []byte) (models.Room, error)
}
