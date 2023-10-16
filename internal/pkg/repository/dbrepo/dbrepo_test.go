package dbrepo

import (
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"homework-3/internal/pkg/models"
	"homework-3/internal/pkg/repository"
	"testing"
)

func Test_postgresDBRepo_GetReservationByID(t *testing.T) {
	t.Parallel()

	var (
		id = 1
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		// arrange

		s := setUp(t)
		defer s.tearDown()
		s.mockDb.EXPECT().Get(gomock.Any(), gomock.Any(), "SELECT id,start_date,end_date,room_id,created_at,updated_at FROM reservations WHERE id=$1", gomock.Any()).Return(nil)

		// act
		reservation, err := s.repo.GetReservationByID(int64(id))
		// assert

		require.NoError(t, err)
		assert.Equal(t, int64(0), reservation.ID)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("not found", func(t *testing.T) {
			t.Parallel()
			// arrange
			s := setUp(t)
			defer s.tearDown()

			s.mockDb.EXPECT().Get(gomock.Any(), gomock.Any(), "SELECT id,start_date,end_date,room_id,created_at,updated_at FROM reservations WHERE id=$1", gomock.Any()).Return(repository.ErrObjectNotFound)

			// act
			reservation, err := s.repo.GetReservationByID(int64(id))
			// assert
			require.EqualError(t, err, "object not found")

			assert.Nil(t, reservation)
		})

		t.Run("internal error", func(t *testing.T) {
			t.Parallel()
			// arrange
			s := setUp(t)
			defer s.tearDown()

			s.mockDb.EXPECT().Get(gomock.Any(), gomock.Any(), "SELECT id,start_date,end_date,room_id,created_at,updated_at FROM reservations WHERE id=$1", gomock.Any()).Return(assert.AnError)

			// act
			reservation, err := s.repo.GetReservationByID(int64(id))
			// assert
			require.EqualError(t, err, "assert.AnError general error for testing")

			assert.Nil(t, reservation)
		})
	})
}

func Test_postgresDBRepo_GetRoomByID(t *testing.T) {
	t.Parallel()

	var (
		id = 1
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		// arrange

		s := setUp(t)
		defer s.tearDown()
		s.mockDb.EXPECT().Get(gomock.Any(), gomock.Any(), "SELECT id,name,cost,created_at,updated_at FROM rooms WHERE id=$1", gomock.Any()).Return(nil)

		// act
		reservation, err := s.repo.GetRoomByID(int64(id))
		// assert

		require.NoError(t, err)
		assert.Equal(t, int64(0), reservation.ID)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("not found", func(t *testing.T) {
			t.Parallel()
			// arrange
			s := setUp(t)
			defer s.tearDown()

			s.mockDb.EXPECT().Get(gomock.Any(), gomock.Any(), "SELECT id,name,cost,created_at,updated_at FROM rooms WHERE id=$1", gomock.Any()).Return(repository.ErrObjectNotFound)

			// act
			reservation, err := s.repo.GetRoomByID(int64(id))
			// assert
			require.EqualError(t, err, "object not found")

			assert.Nil(t, reservation)
		})

		t.Run("internal error", func(t *testing.T) {
			t.Parallel()
			// arrange
			s := setUp(t)
			defer s.tearDown()

			s.mockDb.EXPECT().Get(gomock.Any(), gomock.Any(), "SELECT id,name,cost,created_at,updated_at FROM rooms WHERE id=$1", gomock.Any()).Return(assert.AnError)

			// act
			reservation, err := s.repo.GetRoomByID(int64(id))
			// assert
			require.EqualError(t, err, "assert.AnError general error for testing")

			assert.Nil(t, reservation)
		})
	})
}

//s.mockDb.EXPECT().Exec(gomock.Any(), query, gomock.Any()).Return(pgconn.CommandTag{'1'}, nil)

func Test_postgresDBRepo_GetRoomByName(t *testing.T) {
	t.Parallel()

	var (
		name = "Lux"
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		// arrange

		s := setUp(t)
		defer s.tearDown()
		s.mockDb.EXPECT().Get(gomock.Any(), gomock.Any(), "SELECT id,name,cost,created_at,updated_at FROM rooms WHERE name=$1", gomock.Any()).Return(nil)

		// act
		reservation, err := s.repo.GetRoomByName(name)
		// assert

		require.NoError(t, err)
		assert.Equal(t, int64(0), reservation.ID)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("not found", func(t *testing.T) {
			t.Parallel()
			// arrange
			s := setUp(t)
			defer s.tearDown()

			s.mockDb.EXPECT().Get(gomock.Any(), gomock.Any(), "SELECT id,name,cost,created_at,updated_at FROM rooms WHERE name=$1", gomock.Any()).Return(repository.ErrObjectNotFound)

			// act
			reservation, err := s.repo.GetRoomByName(name)
			// assert
			require.EqualError(t, err, "object not found")

			assert.Nil(t, reservation)
		})

		t.Run("internal error", func(t *testing.T) {
			t.Parallel()
			// arrange
			s := setUp(t)
			defer s.tearDown()

			s.mockDb.EXPECT().Get(gomock.Any(), gomock.Any(), "SELECT id,name,cost,created_at,updated_at FROM rooms WHERE name=$1", gomock.Any()).Return(assert.AnError)

			// act
			reservation, err := s.repo.GetRoomByName(name)
			// assert
			require.EqualError(t, err, "assert.AnError general error for testing")

			assert.Nil(t, reservation)
		})
	})
}

type insertMockRow struct {
	Err error
}

func (im insertMockRow) Scan(dest ...interface{}) error {
	return im.Err
}

func Test_postgresDBRepo_InsertReservation(t *testing.T) {
	t.Parallel()

	var (
		reservation = models.Reservation{}
		query       = `INSERT INTO reservations(start_date, end_date, room_id, created_at, updated_at) VALUES($1,$2,$3,$4,$5) RETURNING id;`
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		// arrange

		s := setUp(t)
		defer s.tearDown()

		row := insertMockRow{nil}
		s.mockDb.EXPECT().ExecQueryRow(gomock.Any(), query, gomock.Any()).Return(row)

		// act
		id, err := s.repo.InsertReservation(&reservation)
		// assert

		require.NoError(t, err)
		assert.Zero(t, id)
	})
	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("internal error", func(t *testing.T) {
			t.Parallel()
			// arrange
			s := setUp(t)
			defer s.tearDown()
			row := insertMockRow{assert.AnError}
			s.mockDb.EXPECT().ExecQueryRow(gomock.Any(), query, gomock.Any()).Return(row)

			// act
			id, err := s.repo.InsertReservation(&reservation)
			// assert
			require.EqualError(t, err, "assert.AnError general error for testing")

			assert.Zero(t, id)
		})
	})
}

func Test_postgresDBRepo_InsertRoom(t *testing.T) {
	t.Parallel()

	var (
		room  = models.Room{}
		query = `INSERT INTO rooms(name, cost, created_at, updated_at) VALUES($1,$2,$3,$4) RETURNING id;`
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		// arrange

		s := setUp(t)
		defer s.tearDown()

		row := insertMockRow{nil}
		s.mockDb.EXPECT().ExecQueryRow(gomock.Any(), query, gomock.Any()).Return(row)

		// act
		id, err := s.repo.InsertRoom(&room)
		// assert

		require.NoError(t, err)
		assert.Equal(t, int64(0), id)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("internal error", func(t *testing.T) {
			t.Parallel()
			// arrange
			s := setUp(t)
			defer s.tearDown()
			row := insertMockRow{assert.AnError}
			s.mockDb.EXPECT().ExecQueryRow(gomock.Any(), query, gomock.Any()).Return(row)

			// act
			id, err := s.repo.InsertRoom(&room)
			// assert
			require.EqualError(t, err, "assert.AnError general error for testing")

			assert.Zero(t, id)
		})
	})
}

func Test_postgresDBRepo_DeleteReservationByID(t *testing.T) {
	t.Parallel()

	var (
		query       = "DELETE FROM reservations WHERE id=$1"
		id    int64 = 1
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		// arrange

		s := setUp(t)
		defer s.tearDown()
		commandTag := pgconn.CommandTag("DELETE 0 1")
		s.mockDb.EXPECT().Exec(gomock.Any(), query, gomock.Any()).Return(commandTag, nil)

		// act
		err := s.repo.DeleteReservationByID(id)
		// assert

		require.NoError(t, err)
	})
	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		// Arrange
		s := setUp(t)
		defer s.tearDown()
		commandTag := pgconn.CommandTag{}
		s.mockDb.EXPECT().Exec(gomock.Any(), query, gomock.Any()).Return(commandTag, nil)

		// Act
		err := s.repo.DeleteReservationByID(id)

		// Assert
		assert.EqualError(t, err, "object not deleted")
	})
}

func Test_postgresDBRepo_DeleteRoomByID(t *testing.T) {
	t.Parallel()

	var (
		query       = "DELETE FROM rooms WHERE id=$1"
		id    int64 = 1
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		// arrange

		s := setUp(t)
		defer s.tearDown()
		commandTag := pgconn.CommandTag("DELETE 0 1")
		s.mockDb.EXPECT().Exec(gomock.Any(), query, gomock.Any()).Return(commandTag, nil)

		// act
		err := s.repo.DeleteRoomByID(id)
		// assert

		require.NoError(t, err)
	})
	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		// Arrange
		s := setUp(t)
		defer s.tearDown()
		commandTag := pgconn.CommandTag{}
		s.mockDb.EXPECT().Exec(gomock.Any(), query, gomock.Any()).Return(commandTag, nil)

		// Act
		err := s.repo.DeleteRoomByID(id)

		// Assert
		assert.EqualError(t, err, "object not deleted")
	})
}
