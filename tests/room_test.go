//go:build integration
// +build integration

package tests

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"homework-3/internal/pkg/repository/dbrepo"
	"homework-3/tests/fixtures"
	"testing"
)

func Test_InsertRoom(t *testing.T) {
	t.Parallel()
	var room = fixtures.Room().Valid().P()
	ctx := context.Background()
	t.Run("success", func(t *testing.T) {
		db.SetUp(t)
		defer db.TearDown()
		//arrange

		repo := dbrepo.NewPostgresRepo(db.DB)

		//act
		resp, err := repo.InsertRoom(ctx, room)

		//assert
		require.NoError(t, err)
		assert.NotZero(t, resp)
	})

	t.Run("fail - repeated name", func(t *testing.T) {
		db.SetUp(t)
		defer db.TearDown()
		//arrange

		repo := dbrepo.NewPostgresRepo(db.DB)
		resp1, err := repo.InsertRoom(ctx, room)

		require.NoError(t, err)
		assert.NotZero(t, resp1)

		//act
		resp2, err := repo.InsertRoom(ctx, room)

		//assert
		require.EqualError(t, err, "ERROR: duplicate key value violates unique constraint \"rooms_name_key\" (SQLSTATE 23505)")
		assert.Zero(t, resp2)
	})
}

func Test_UpdateRoom(t *testing.T) {
	var room = fixtures.Room().Valid().P()
	var updateRoom = fixtures.Room().UpdatedValidWithDifferentCost().P()
	var ctx = context.Background()
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		roomID, err := repo.InsertRoom(ctx, room)

		require.NoError(t, err)
		assert.NotZero(t, roomID)

		roomBefore, err := repo.GetRoomByID(ctx, roomID)
		require.NoError(t, err)

		updateRoom.ID = roomID

		//act
		err = repo.UpdateRoom(ctx, updateRoom)

		//assert

		require.NoError(t, err)

		roomAfter, err := repo.GetRoomByID(ctx, roomID)
		require.NoError(t, err)
		assert.NotEqual(t, roomBefore.Cost, roomAfter.Cost)
	})

	t.Run("fail - invalid id", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		roomID, err := repo.InsertRoom(ctx, room)

		require.NoError(t, err)
		assert.NotZero(t, roomID)

		roomBefore, err := repo.GetRoomByID(ctx, roomID)
		require.NoError(t, err)

		updateRoom.ID = roomID + 1

		//act
		err = repo.UpdateRoom(ctx, updateRoom)

		//assert
		require.EqualError(t, err, "object not updated")

		roomAfter, err := repo.GetRoomByID(ctx, roomID)
		require.NoError(t, err)
		assert.Equal(t, roomBefore.Cost, roomAfter.Cost)
	})
}

func Test_GetRoomByID(t *testing.T) {
	var (
		room = fixtures.Room().Valid().P()
		ctx  = context.Background()
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		resp, err := repo.InsertRoom(ctx, room)

		require.NoError(t, err)
		assert.NotZero(t, resp)

		//act
		getRoom, err := repo.GetRoomByID(ctx, resp)

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
		resp, err := repo.InsertRoom(ctx, room)

		require.NoError(t, err)
		assert.NotZero(t, resp)

		//act
		getRoom, err := repo.GetRoomByID(ctx, resp+1)

		//assert
		require.EqualError(t, err, "object not found")
		assert.Nil(t, getRoom)
	})
}

func Test_GetRoomByName(t *testing.T) {
	var (
		room = fixtures.Room().Valid().P()
		ctx  = context.Background()
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		resp, err := repo.InsertRoom(ctx, room)

		require.NoError(t, err)
		assert.NotZero(t, resp)

		//act
		getRoom, err := repo.GetRoomByName(ctx, room.Name)

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
		resp, err := repo.InsertRoom(ctx, room)

		require.NoError(t, err)
		assert.NotZero(t, resp)

		//act
		getRoom, err := repo.GetRoomByName(ctx, room.Name+"_invalid_name")

		//assert
		require.EqualError(t, err, "object not found")
		assert.Nil(t, getRoom)
	})
}

func Test_DeleteRoomByID(t *testing.T) {
	var (
		room = fixtures.Room().Valid().P()
		ctx  = context.Background()
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		roomID, err := repo.InsertRoom(ctx, room)

		require.NoError(t, err)
		assert.NotZero(t, roomID)

		//act
		err = repo.DeleteRoomByID(ctx, roomID)

		//assert
		require.NoError(t, err)

		getRoom, err := repo.GetRoomByID(ctx, roomID)
		require.Error(t, err, "object not found")
		require.Nil(t, getRoom)
	})

	t.Run("fail - invalid id", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		roomID, err := repo.InsertRoom(ctx, room)

		require.NoError(t, err)
		assert.NotZero(t, roomID)

		//act
		err = repo.DeleteRoomByID(ctx, roomID+1)

		//assert
		require.EqualError(t, err, "object not deleted")

		getRoom, err := repo.GetRoomByID(ctx, roomID)
		require.Nil(t, err)
		assert.Equal(t, room.Name, getRoom.Name)
		assert.Equal(t, room.Cost, getRoom.Cost)
		assert.Equal(t, room.CreatedAt, getRoom.CreatedAt)
		assert.Equal(t, room.UpdatedAt, getRoom.UpdatedAt)
	})
}

func Test_InsertReservation(t *testing.T) {
	t.Parallel()
	var (
		res = fixtures.Reservation().Valid().P()
		ctx = context.Background()
	)
	t.Run("success", func(t *testing.T) {
		db.SetUp(t)
		defer db.TearDown()
		//arrange

		repo := dbrepo.NewPostgresRepo(db.DB)

		//act
		resp, err := repo.InsertReservation(ctx, res)

		//assert
		require.NoError(t, err)
		assert.NotZero(t, resp)
	})
}

func Test_UpdateReservation(t *testing.T) {
	var (
		res       = fixtures.Reservation().Valid().P()
		updateRes = fixtures.Reservation().UpdatedValidWithDifferentDates().P()
		ctx       = context.Background()
	)
	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		resID, err := repo.InsertReservation(ctx, res)

		require.NoError(t, err)
		assert.NotZero(t, resID)

		resBefore, err := repo.GetReservationByID(ctx, resID)
		require.NoError(t, err)

		updateRes.ID = resID

		//act
		err = repo.UpdateReservation(ctx, updateRes)

		//assert
		require.NoError(t, err)

		resAfter, err := repo.GetReservationByID(ctx, resID)
		require.NoError(t, err)

		require.NotEqual(t, resBefore.StartDate, resAfter.StartDate)
		require.NotEqual(t, resBefore.EndDate, resAfter.EndDate)
	})

	t.Run("fail - invalid id", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		resID, err := repo.InsertReservation(ctx, res)

		require.NoError(t, err)
		assert.NotZero(t, resID)

		resBefore, err := repo.GetReservationByID(ctx, resID)
		require.NoError(t, err)

		updateRes.ID = resID + 1

		//act
		err = repo.UpdateReservation(ctx, updateRes)

		//assert
		require.EqualError(t, err, "object not updated")

		resAfter, err := repo.GetReservationByID(ctx, resID)
		require.NoError(t, err)

		assert.Equal(t, resBefore.StartDate, resAfter.StartDate)
		assert.Equal(t, resBefore.EndDate, resAfter.EndDate)
	})
}

func Test_GetReservationByID(t *testing.T) {
	var (
		res = fixtures.Reservation().Valid().P()
		ctx = context.Background()
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		resp, err := repo.InsertReservation(ctx, res)

		require.NoError(t, err)
		assert.NotZero(t, resp)

		//act
		getRes, err := repo.GetReservationByID(ctx, resp)

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
		resp, err := repo.InsertReservation(ctx, res)

		require.NoError(t, err)
		assert.NotZero(t, resp)

		//act
		getRoom, err := repo.GetReservationByID(ctx, resp+1)

		//assert
		require.EqualError(t, err, "object not found")
		assert.Nil(t, getRoom)
	})
}

func Test_DeleteReservationByID(t *testing.T) {
	var (
		res = fixtures.Reservation().Valid().P()
		ctx = context.Background()
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		resID, err := repo.InsertReservation(ctx, res)

		require.NoError(t, err)
		assert.NotZero(t, resID)

		//act
		err = repo.DeleteReservationByID(ctx, resID)

		//assert
		require.NoError(t, err)

		getRes, err := repo.GetReservationByID(ctx, resID)
		require.Error(t, err, "object not found")
		require.Nil(t, getRes)
	})

	t.Run("fail - invalid id", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		resID, err := repo.InsertReservation(ctx, res)

		require.NoError(t, err)
		assert.NotZero(t, resID)

		//act
		err = repo.DeleteReservationByID(ctx, resID+1)

		//assert
		require.EqualError(t, err, "object not deleted")

		getRes, err := repo.GetReservationByID(ctx, resID)
		require.Nil(t, err)
		require.Equal(t, res.RoomID, getRes.RoomID)
		require.Equal(t, res.StartDate, getRes.StartDate)
		require.Equal(t, res.EndDate, getRes.EndDate)
		require.Equal(t, res.CreatedAt, getRes.CreatedAt)
		require.Equal(t, res.UpdatedAt, getRes.UpdatedAt)
	})
}

func Test_GetReservationsByRoomID(t *testing.T) {
	var (
		room = fixtures.Room().Valid().P()
		res  = fixtures.Reservation().Valid().P()
		res2 = fixtures.Reservation().Valid2().P()
		ctx  = context.Background()
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		roomResp, err := repo.InsertRoom(ctx, room)
		require.NoError(t, err)
		require.NotZero(t, roomResp)

		res.RoomID = roomResp

		resResp1, err := repo.InsertReservation(ctx, res)
		require.NoError(t, err)
		require.NotZero(t, resResp1)

		res2.RoomID = roomResp

		resResp2, err := repo.InsertReservation(ctx, res2)
		require.NoError(t, err)
		require.NotZero(t, resResp2)

		//act
		reservations, err := repo.GetReservationsByRoomID(ctx, roomResp)

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
		roomResp, err := repo.InsertRoom(ctx, room)
		require.NoError(t, err)
		require.NotZero(t, roomResp)

		res.RoomID = roomResp

		resResp1, err := repo.InsertReservation(ctx, res)
		require.NoError(t, err)
		require.NotZero(t, resResp1)

		res2.RoomID = roomResp

		resResp2, err := repo.InsertReservation(ctx, res2)
		require.NoError(t, err)
		require.NotZero(t, resResp2)

		//act
		reservations, err := repo.GetReservationsByRoomID(ctx, roomResp+1)

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
		ctx  = context.Background()
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		roomID, err := repo.InsertRoom(ctx, room)
		require.NoError(t, err)
		require.NotZero(t, roomID)

		res.RoomID = roomID

		resResp1, err := repo.InsertReservation(ctx, res)
		require.NoError(t, err)
		require.NotZero(t, resResp1)

		res2.RoomID = roomID

		resResp2, err := repo.InsertReservation(ctx, res2)
		require.NoError(t, err)
		require.NotZero(t, resResp2)

		//act
		err = repo.DeleteReservationsByRoomID(ctx, roomID)

		//assert
		require.NoError(t, err)

		getReservations, err := repo.GetReservationsByRoomID(ctx, roomID)

		require.Nil(t, getReservations)
		require.Nil(t, err)

	})

	t.Run("fail - invalid id", func(t *testing.T) {
		t.Parallel()
		db.SetUp(t)
		defer db.TearDown()
		//arrange
		repo := dbrepo.NewPostgresRepo(db.DB)
		roomID, err := repo.InsertRoom(ctx, room)
		require.NoError(t, err)
		require.NotZero(t, roomID)

		res.RoomID = roomID

		resResp1, err := repo.InsertReservation(ctx, res)
		require.NoError(t, err)
		require.NotZero(t, resResp1)

		res2.RoomID = roomID

		resResp2, err := repo.InsertReservation(ctx, res2)
		require.NoError(t, err)
		require.NotZero(t, resResp2)

		//act
		err = repo.DeleteReservationsByRoomID(ctx, roomID+1)

		//assert
		require.EqualError(t, err, "object not found")

		getReservations, err := repo.GetReservationsByRoomID(ctx, roomID)

		require.NoError(t, err)

		assert.Equal(t, res.StartDate, getReservations[0].StartDate)
		assert.Equal(t, res.EndDate, getReservations[0].EndDate)
		assert.Equal(t, res.RoomID, getReservations[0].RoomID)

		assert.Equal(t, res2.StartDate, getReservations[1].StartDate)
		assert.Equal(t, res2.EndDate, getReservations[1].EndDate)
		assert.Equal(t, res2.RoomID, getReservations[1].RoomID)
	})
}
