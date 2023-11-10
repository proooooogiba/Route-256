package grpc_server

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	"homework-3/internal/grpc_handlers"
	"homework-3/internal/pkg/models"
	"homework-3/internal/pkg/pb"
	"time"
)

type Transport struct {
	pb.UnimplementedHotelServiceServer
	hotelService *grpc_handlers.Service
}

func NewImplementation(service *grpc_handlers.Service) *Transport {
	return &Transport{
		hotelService: service,
	}
}

func (t *Transport) GetRoomWithAllReservations(ctx context.Context, in *pb.GetRoomRequest) (*pb.GetRoomResponse, error) {
	room, reservations, err := t.hotelService.GetRoomWithAllReservations(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	layout := "2006-01-02"

	getReservationsResponse := make([]*pb.GetReservationResponse, len(reservations))

	for i, res := range reservations {
		getRes := &pb.GetReservationResponse{
			Id:        res.ID,
			StartDate: res.StartDate.Format(layout),
			EndDate:   res.EndDate.Format(layout),
			RoomId:    res.RoomID,
			CreatedAt: timestamppb.New(res.CreatedAt),
			UpdatedAt: timestamppb.New(res.UpdatedAt),
		}
		getReservationsResponse[i] = getRes
	}

	return &pb.GetRoomResponse{
		Id:           room.ID,
		Name:         room.Name,
		Cost:         room.Cost,
		CreatedAt:    timestamppb.New(room.CreatedAt),
		UpdatedAt:    timestamppb.New(room.UpdatedAt),
		Reservations: getReservationsResponse,
	}, nil
}

func (t *Transport) CreateRoom(ctx context.Context, in *pb.CreateRoomRequest) (*pb.CreateRoomResponse, error) {
	room := models.Room{
		Name:      in.Name,
		Cost:      in.Cost,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	roomID, err := t.hotelService.CreateRoom(ctx, room)
	if err != nil {
		return nil, err
	}

	return &pb.CreateRoomResponse{
		Id: roomID,
	}, nil
}

func (t *Transport) UpdateRoom(ctx context.Context, in *pb.UpdateRoomRequest) (*pb.UpdateRoomResponse, error) {
	room := models.Room{
		ID:        in.Id,
		Name:      in.Name,
		Cost:      in.Cost,
		UpdatedAt: time.Now(),
	}

	err := t.hotelService.UpdateRoom(ctx, room)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateRoomResponse{
		Ok: true,
	}, nil
}

func (t *Transport) DeleteRoomWithAllReservations(ctx context.Context, in *pb.DeleteRoomRequest) (*pb.DeleteRoomResponse, error) {
	err := t.hotelService.DeleteRoomWithAllReservations(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteRoomResponse{
		Ok: true,
	}, nil
}

func (t *Transport) CreateReservation(ctx context.Context, in *pb.CreateReservationRequest) (*pb.CreateReservationResponse, error) {
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, in.StartDate)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse(layout, in.EndDate)
	if err != nil {
		return nil, err
	}

	reservation := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    in.RoomId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	resID, err := t.hotelService.CreateReservation(ctx, reservation)
	if err != nil {
		return nil, err
	}

	return &pb.CreateReservationResponse{
		Id: resID,
	}, nil
}

func (t *Transport) GetReservation(ctx context.Context, in *pb.GetReservationRequest) (*pb.GetReservationResponse, error) {
	room, err := t.hotelService.GetReservation(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	layout := "2006-01-02"
	return &pb.GetReservationResponse{
		Id:        room.ID,
		StartDate: room.StartDate.Format(layout),
		EndDate:   room.EndDate.Format(layout),
		RoomId:    room.RoomID,
		CreatedAt: timestamppb.New(room.CreatedAt),
		UpdatedAt: timestamppb.New(room.UpdatedAt),
	}, nil
}

func (t *Transport) UpdateReservation(ctx context.Context, in *pb.UpdateReservationRequest) (*pb.UpdateReservationResponse, error) {
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, in.StartDate)
	if err != nil {
		return nil, err
	}

	endDate, err := time.Parse(layout, in.EndDate)
	if err != nil {
		return nil, err
	}

	reservation := models.Reservation{
		ID:        in.Id,
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    in.RoomId,
		UpdatedAt: time.Now(),
	}

	err = t.hotelService.UpdateReservation(ctx, reservation)
	if err != nil {
		return nil, err
	}

	return &pb.UpdateReservationResponse{
		Ok: true,
	}, nil
}

func (t *Transport) DeleteReservation(ctx context.Context, in *pb.DeleteReservationRequest) (*pb.DeleteReservationResponse, error) {
	err := t.hotelService.DeleteReservation(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	return &pb.DeleteReservationResponse{
		Ok: true,
	}, nil
}
