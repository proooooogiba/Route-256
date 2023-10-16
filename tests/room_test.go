//go:build integration
// +build integration

package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"homework-3/internal/pkg/repository/dbrepo"
	"homework-3/tests/fixtures"
	"testing"
)

func Test_InsertRoom(t *testing.T) {
	t.Parallel()
	var room = fixtures.Room().Valid().P()
	t.Run("success", func(t *testing.T) {
		db.SetUp(t)
		defer db.TearDown()
		//arrange

		repo := dbrepo.NewPostgresRepo(db.DB)

		//act
		resp, err := repo.InsertRoom(room)

		//assert
		require.NoError(t, err)
		assert.NotZero(t, resp)
	})

	t.Run("fail - repeated name", func(t *testing.T) {
		db.SetUp(t)
		defer db.TearDown()
		//arrange

		repo := dbrepo.NewPostgresRepo(db.DB)
		resp1, err := repo.InsertRoom(room)

		require.NoError(t, err)
		assert.NotZero(t, resp1)

		//act
		resp2, err := repo.InsertRoom(room)

		//assert
		require.EqualError(t, err, "ERROR: duplicate key value violates unique constraint \"rooms_name_key\" (SQLSTATE 23505)")
		assert.Zero(t, resp2)
	})
}

func Test_UpdateRoom(t *testing.T) {
	var room = fixtures.Room().Valid().P()
	var update_room = fixtures.Room().UpdatedValid().P()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		resp, err := repo.InsertRoom(room)

		require.NoError(t, err)
		assert.NotZero(t, resp)

		update_room.ID = resp

		//act
		err = repo.UpdateRoom(update_room)

		//assert
		require.NoError(t, err)
	})

	t.Run("fail - invalid id", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		resp, err := repo.InsertRoom(room)

		require.NoError(t, err)
		assert.NotZero(t, resp)

		update_room.ID = resp + 1

		//act
		err = repo.UpdateRoom(update_room)

		//assert
		require.EqualError(t, err, "object not updated")
	})
}

func Test_GetRoomByID(t *testing.T) {
	var (
		room = fixtures.Room().Valid().P()
		//update_room = fixtures.Room().UpdatedValid().P()
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		resp, err := repo.InsertRoom(room)

		require.NoError(t, err)
		assert.NotZero(t, resp)

		//act
		getRoom, err := repo.GetRoomByID(resp)

		//assert
		require.NoError(t, err)
		assert.Equal(t, room.Name, getRoom.Name)
		assert.Equal(t, room.Cost, getRoom.Cost)
		assert.Equal(t, room.CreatedAt, getRoom.CreatedAt)
		assert.Equal(t, room.UpdatedAt, getRoom.UpdatedAt)
	})

	t.Run("fail - invalid id", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		resp, err := repo.InsertRoom(room)

		require.NoError(t, err)
		assert.NotZero(t, resp)

		//act
		getRoom, err := repo.GetRoomByID(resp + 1)

		//assert
		require.EqualError(t, err, "object not found")
		assert.Nil(t, getRoom)
	})
}

func Test_GetRoomByName(t *testing.T) {
	var (
		room = fixtures.Room().Valid().P()
		//update_room = fixtures.Room().UpdatedValid().P()
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		resp, err := repo.InsertRoom(room)

		require.NoError(t, err)
		assert.NotZero(t, resp)

		//act
		getRoom, err := repo.GetRoomByName(room.Name)

		//assert
		require.NoError(t, err)
		assert.Equal(t, room.Name, getRoom.Name)
		assert.Equal(t, room.Cost, getRoom.Cost)
		assert.Equal(t, room.CreatedAt, getRoom.CreatedAt)
		assert.Equal(t, room.UpdatedAt, getRoom.UpdatedAt)
	})

	t.Run("fail - invalid name", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		resp, err := repo.InsertRoom(room)

		require.NoError(t, err)
		assert.NotZero(t, resp)

		//act
		getRoom, err := repo.GetRoomByName(room.Name + "_invalid_name")

		//assert
		require.EqualError(t, err, "object not found")
		assert.Nil(t, getRoom)
	})
}

func Test_DeleteRoomByID(t *testing.T) {
	var (
		room = fixtures.Room().Valid().P()
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		resp, err := repo.InsertRoom(room)

		require.NoError(t, err)
		assert.NotZero(t, resp)

		//act
		err = repo.DeleteRoomByID(resp)

		//assert
		require.NoError(t, err)
	})

	t.Run("fail - invalid id", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		resp, err := repo.InsertRoom(room)

		require.NoError(t, err)
		assert.NotZero(t, resp)

		//act
		err = repo.DeleteRoomByID(resp + 1)

		//assert
		require.EqualError(t, err, "object not deleted")
	})
}

func Test_InsertReservation(t *testing.T) {
	t.Parallel()
	var res = fixtures.Reservation().Valid().P()
	t.Run("success", func(t *testing.T) {
		db.SetUp(t)
		defer db.TearDown()
		//arrange

		repo := dbrepo.NewPostgresRepo(db.DB)

		//act
		resp, err := repo.InsertReservation(res)

		//assert
		require.NoError(t, err)
		assert.NotZero(t, resp)
	})
}

func Test_UpdateReservation(t *testing.T) {
	var res = fixtures.Reservation().Valid().P()
	var updateRes = fixtures.Reservation().UpdatedValid().P()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		resp, err := repo.InsertReservation(res)

		require.NoError(t, err)
		assert.NotZero(t, resp)

		updateRes.ID = resp

		//act
		err = repo.UpdateReservation(updateRes)

		//assert
		require.NoError(t, err)
	})

	t.Run("fail - invalid id", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		resp, err := repo.InsertReservation(res)

		require.NoError(t, err)
		assert.NotZero(t, resp)

		updateRes.ID = resp + 1

		//act
		err = repo.UpdateReservation(updateRes)

		//assert
		require.EqualError(t, err, "object not updated")
	})
}

func Test_GetReservationByID(t *testing.T) {
	var (
		res = fixtures.Reservation().Valid().P()
		//update_room = fixtures.Room().UpdatedValid().P()
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		resp, err := repo.InsertReservation(res)

		require.NoError(t, err)
		assert.NotZero(t, resp)

		//act
		getRes, err := repo.GetReservationByID(resp)

		//assert
		require.NoError(t, err)
		assert.Equal(t, res.StartDate, getRes.StartDate)
		assert.Equal(t, res.EndDate, getRes.EndDate)
		assert.Equal(t, res.RoomID, getRes.RoomID)
		assert.Equal(t, res.CreatedAt, getRes.CreatedAt)
		assert.Equal(t, res.UpdatedAt, getRes.UpdatedAt)
	})

	t.Run("fail - invalid id", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		resp, err := repo.InsertReservation(res)

		require.NoError(t, err)
		assert.NotZero(t, resp)

		//act
		getRoom, err := repo.GetReservationByID(resp + 1)

		//assert
		require.EqualError(t, err, "object not found")
		assert.Nil(t, getRoom)
	})
}

func Test_DeleteReservationByID(t *testing.T) {
	var (
		res = fixtures.Reservation().Valid().P()
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		resp, err := repo.InsertReservation(res)

		require.NoError(t, err)
		assert.NotZero(t, resp)

		//act
		err = repo.DeleteReservationByID(resp)

		//assert
		require.NoError(t, err)
	})

	t.Run("fail - invalid id", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		resp, err := repo.InsertReservation(res)

		require.NoError(t, err)
		assert.NotZero(t, resp)

		//act
		err = repo.DeleteReservationByID(resp + 1)

		//assert
		require.EqualError(t, err, "object not deleted")
	})
}

func Test_GetReservationsByRoomID(t *testing.T) {
	var (
		room = fixtures.Room().Valid().P()
		res  = fixtures.Reservation().Valid().P()
		res2 = fixtures.Reservation().Valid2().P()
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		roomResp, err := repo.InsertRoom(room)
		require.NoError(t, err)
		require.NotZero(t, roomResp)

		res.RoomID = roomResp

		resResp1, err := repo.InsertReservation(res)
		require.NoError(t, err)
		require.NotZero(t, resResp1)

		res2.RoomID = roomResp

		resResp2, err := repo.InsertReservation(res2)
		require.NoError(t, err)
		require.NotZero(t, resResp2)

		//act
		reservations, err := repo.GetReservationsByRoomID(roomResp)

		//assert
		require.NoError(t, err)
		assert.Equal(t, res.StartDate, reservations[0].StartDate)
		assert.Equal(t, res.EndDate, reservations[0].EndDate)
		assert.Equal(t, res.RoomID, reservations[0].RoomID)

		assert.Equal(t, res2.StartDate, reservations[1].StartDate)
		assert.Equal(t, res2.EndDate, reservations[1].EndDate)
		assert.Equal(t, res2.RoomID, reservations[1].RoomID)
	})

	t.Run("fail - invalid id", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		roomResp, err := repo.InsertRoom(room)
		require.NoError(t, err)
		require.NotZero(t, roomResp)

		res.RoomID = roomResp

		resResp1, err := repo.InsertReservation(res)
		require.NoError(t, err)
		require.NotZero(t, resResp1)

		res2.RoomID = roomResp

		resResp2, err := repo.InsertReservation(res2)
		require.NoError(t, err)
		require.NotZero(t, resResp2)

		//act
		reservations, err := repo.GetReservationsByRoomID(roomResp + 1)

		//assert
		require.EqualError(t, err, "object not found")
		assert.Nil(t, reservations)
	})
}

func Test_DeleteReservationsByRoomID(t *testing.T) {
	var (
		room = fixtures.Room().Valid().P()
		res  = fixtures.Reservation().Valid().P()
		res2 = fixtures.Reservation().Valid2().P()
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		roomResp, err := repo.InsertRoom(room)
		require.NoError(t, err)
		require.NotZero(t, roomResp)

		res.RoomID = roomResp

		resResp1, err := repo.InsertReservation(res)
		require.NoError(t, err)
		require.NotZero(t, resResp1)

		res2.RoomID = roomResp

		resResp2, err := repo.InsertReservation(res2)
		require.NoError(t, err)
		require.NotZero(t, resResp2)

		//act
		err = repo.DeleteReservationsByRoomID(roomResp)

		//assert
		require.NoError(t, err)
	})

	t.Run("fail - invalid id", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		roomResp, err := repo.InsertRoom(room)
		require.NoError(t, err)
		require.NotZero(t, roomResp)

		res.RoomID = roomResp

		resResp1, err := repo.InsertReservation(res)
		require.NoError(t, err)
		require.NotZero(t, resResp1)

		res2.RoomID = roomResp

		resResp2, err := repo.InsertReservation(res2)
		require.NoError(t, err)
		require.NotZero(t, resResp2)

		//act
		err = repo.DeleteReservationsByRoomID(roomResp + 1)

		//assert
		require.EqualError(t, err, "object not found")
	})
}
