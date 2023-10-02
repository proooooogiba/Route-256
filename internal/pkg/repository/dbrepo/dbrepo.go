package dbrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"homework-3/internal/pkg/db"
	"homework-3/internal/pkg/models"
	"homework-3/internal/pkg/repository"
)

type postgresDBRepo struct {
	db *db.Database
}

func NewPostgresRepo(db *db.Database) repository.DatabaseRepo {
	return &postgresDBRepo{
		db: db,
	}
}

func (r *postgresDBRepo) InsertReservation(reservation models.Reservation) (int64, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var newID int64
	err := r.db.ExecQueryRow(ctx, `INSERT INTO reservations(start_date, end_date, room_id, created_at, updated_at) VALUES($1,$2,$3,$4,$5) RETURNING id;`, reservation.StartDate, reservation.EndDate, reservation.RoomID, reservation.CreatedAt, reservation.UpdatedAt).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (r *postgresDBRepo) GetReservationByID(id int64) (models.Reservation, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var res models.Reservation

	err := r.db.Get(ctx, &res, "SELECT id,start_date,end_date,room_id,created_at,updated_at FROM reservations WHERE id=$1", id)
	if errors.Is(err, sql.ErrNoRows) {
		return res, repository.ErrObjectNotFound
	}
	return res, nil
}

func (r *postgresDBRepo) DeleteReservationByID(id int64) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result, err := r.db.Exec(ctx, "DELETE FROM reservations WHERE id=$1", id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return repository.ErrObjectNotDelete
	}

	return nil
}

func (r *postgresDBRepo) UpdateReservation(res models.Reservation) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fmt.Println(res)
	result, err := r.db.Exec(ctx, "UPDATE reservations set start_date = $1, end_date = $2, room_id = $3, updated_at = $4 where id = $5",
		res.StartDate,
		res.EndDate,
		res.RoomID,
		res.UpdatedAt,
		res.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return repository.ErrObjectNotUpdated
	}
	return nil
}

func (r *postgresDBRepo) InsertRoom(room models.Room) (int64, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var newID int64
	err := r.db.ExecQueryRow(ctx, "INSERT INTO rooms(name, cost, created_at, updated_at) VALUES($1,$2,$3,$4) RETURNING id;", room.Name, room.Cost, room.CreatedAt, room.UpdatedAt).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil
}

func (r *postgresDBRepo) GetRoomByID(id int64) (models.Room, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var room models.Room

	err := r.db.Get(ctx, &room, "SELECT id,name,cost,created_at,updated_at FROM rooms WHERE id=$1", id)

	if errors.Is(err, pgx.ErrNoRows) {
		return room, repository.ErrObjectNotFound
	} else if err != nil {
		return room, err
	}
	return room, nil
}

func (r *postgresDBRepo) GetRoomByName(name string) (models.Room, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var room models.Room

	err := r.db.Get(ctx, &room, "SELECT id,name,cost,created_at,updated_at FROM rooms WHERE name=$1", name)

	if errors.Is(err, pgx.ErrNoRows) {
		return room, repository.ErrObjectNotFound
	} else if err != nil {
		return room, err
	}
	return room, nil
}

func (r *postgresDBRepo) DeleteRoomByID(id int64) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result, err := r.db.Exec(ctx, "DELETE FROM rooms WHERE id=$1", id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return repository.ErrObjectNotDelete
	}

	return nil
}

func (r *postgresDBRepo) UpdateRoom(room models.Room) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	result, err := r.db.Exec(ctx, "UPDATE rooms set name = $1, cost = $2, updated_at = $3 where id = $4",
		room.Name,
		room.Cost,
		room.UpdatedAt,
		room.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return repository.ErrObjectNotUpdated
	}
	return nil
}

func (r *postgresDBRepo) GetReservationsByRoomID(roomID int64) ([]models.Reservation, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var reservations []models.Reservation

	rows, err := r.db.ExecQueryRows(ctx, "SELECT id,start_date,end_date,room_id,created_at,updated_at FROM reservations WHERE room_id=$1", roomID)
	if err != nil {
		return reservations, err
	}

	for rows.Next() {
		var res models.Reservation
		err := rows.Scan(
			&res.ID,
			&res.StartDate,
			&res.EndDate,
			&res.RoomID,
			&res.CreatedAt,
			&res.UpdatedAt,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, res)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	return reservations, nil
}

func (r *postgresDBRepo) DeleteReservationsByRoomID(roomID int64) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := r.db.Exec(ctx, "DELETE FROM reservations WHERE room_id=$1", roomID)
	if err != nil {
		return err
	}
	return nil
}
