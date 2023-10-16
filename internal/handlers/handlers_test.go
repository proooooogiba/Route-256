package handlers

import (
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"homework-3/internal/pkg/models"
	"homework-3/internal/pkg/repository"
	mock_repository "homework-3/internal/pkg/repository/mocks"
	"homework-3/tests/fixtures"
	"net/http"
	"testing"
)

func Test_GetRoomWithAllReservations(t *testing.T) {
	var (
		id   int64 = 1
		room       = fixtures.Room().Valid().P()
	)

	reservations := make([]*models.Reservation, 2, 4)
	reservations = append(reservations, fixtures.Reservation().Valid().P())
	reservations = append(reservations, fixtures.Reservation().Valid2().P())

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		t.Run("without reservations", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_repository.NewMockDatabaseRepo(ctrl)
			s := NewRepo(m)

			m.EXPECT().GetRoomByID(id).Return(room, nil)
			m.EXPECT().GetReservationsByRoomID(id).Return(nil, nil)

			//act
			getRoom, getReservations, code := s.GetRoomWithAllReservations(id)

			// assert
			require.Equal(t, http.StatusOK, code)
			assert.Equal(t, "{\"id\":1,\"name\":\"Lux\",\"cost\":1000,\"created_at\":\"0001-01-01T00:00:00Z\",\"updated_at\":\"0001-01-01T00:00:00Z\"}", string(getRoom))
			assert.Nil(t, getReservations)
		})
		t.Run("with reservations", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_repository.NewMockDatabaseRepo(ctrl)
			s := NewRepo(m)

			m.EXPECT().GetRoomByID(id).Return(room, nil)
			m.EXPECT().GetReservationsByRoomID(id).Return(reservations, nil)
			//act
			getRoom, getReservations, code := s.GetRoomWithAllReservations(id)

			// assert
			require.Equal(t, http.StatusOK, code)
			assert.Equal(t, "{\"id\":1,\"name\":\"Lux\",\"cost\":1000,\"created_at\":\"0001-01-01T00:00:00Z\",\"updated_at\":\"0001-01-01T00:00:00Z\"}", string(getRoom))
			assert.Equal(t, "\nnull\nnull\n{\"id\":1,\"start_date\":\"2023-11-08T00:00:00Z\",\"end_date\":\"2023-11-17T00:00:00Z\",\"room_id\":1,\"created_at\":\"0001-01-01T00:00:00Z\",\"updated_at\":\"0001-01-01T00:00:00Z\"}\n"+
				"{\"id\":2,\"start_date\":\"2023-10-08T00:00:00Z\",\"end_date\":\"2023-10-17T00:00:00Z\",\"room_id\":1,\"created_at\":\"0001-01-01T00:00:00Z\",\"updated_at\":\"0001-01-01T00:00:00Z\"}\n", string(getReservations))
		})
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()

		t.Run("room not found", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_repository.NewMockDatabaseRepo(ctrl)
			s := NewRepo(m)

			m.EXPECT().GetRoomByID(id).Return(nil, repository.ErrObjectNotFound)
			//act
			getRoom, getReservations, code := s.GetRoomWithAllReservations(id)

			// assert
			require.Equal(t, http.StatusNotFound, code)
			assert.Nil(t, getRoom, getReservations)
		})

		t.Run("internal error", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_repository.NewMockDatabaseRepo(ctrl)
			s := NewRepo(m)

			m.EXPECT().GetRoomByID(id).Return(room, nil)
			m.EXPECT().GetReservationsByRoomID(id).Return(nil, errors.New("Error while getting reservations by room id"))
			//act
			getRoom, getReservations, code := s.GetRoomWithAllReservations(id)

			// assert
			require.Equal(t, http.StatusInternalServerError, code)
			assert.Nil(t, getReservations, getRoom)
		})
	})
}

func Test_CreateRoom(t *testing.T) {
	var (
		id      int64 = 1
		name          = "Lux"
		roomReq       = "{\"name\":\"Lux\", \"cost\":1000.0}"
	)
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := mock_repository.NewMockDatabaseRepo(ctrl)
		s := NewRepo(m)

		m.EXPECT().GetRoomByName(name).Return(nil, repository.ErrObjectNotFound)
		m.EXPECT().InsertRoom(gomock.Any()).Return(id, nil)
		//act
		code := s.CreateRoom([]byte(roomReq))

		// assert
		require.Equal(t, http.StatusOK, code)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("room with following name already exists", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_repository.NewMockDatabaseRepo(ctrl)
			s := NewRepo(m)

			m.EXPECT().GetRoomByName(name).Return(nil, nil)
			//act
			code := s.CreateRoom([]byte(roomReq))

			// assert
			require.Equal(t, http.StatusConflict, code)
		})
		t.Run("room with following name already exists", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_repository.NewMockDatabaseRepo(ctrl)
			s := NewRepo(m)

			m.EXPECT().GetRoomByName(name).Return(nil, repository.ErrObjectNotFound)
			m.EXPECT().InsertRoom(gomock.Any()).Return(id, repository.ErrInternalServer)

			//act
			code := s.CreateRoom([]byte(roomReq))

			// assert
			require.Equal(t, http.StatusInternalServerError, code)
		})
	})
}

func Test_UpdateRoom(t *testing.T) {
	var (
		id int64 = 1
		//name          = "Lux"
		roomReq = "{\"id\":1, \"name\":\"Lux\", \"cost\":1000.0}"
		room    = models.Room{
			ID:   1,
			Name: "Lux",
			Cost: 1000.0,
		}
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := mock_repository.NewMockDatabaseRepo(ctrl)
		s := NewRepo(m)

		m.EXPECT().GetRoomByID(id).Return(&room, nil)

		m.EXPECT().UpdateRoom(gomock.Any()).Return(nil)
		//act
		code := s.UpdateRoom([]byte(roomReq))

		// assert
		require.Equal(t, http.StatusOK, code)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("room with id doesn't exists", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_repository.NewMockDatabaseRepo(ctrl)
			s := NewRepo(m)

			m.EXPECT().GetRoomByID(id).Return(nil, repository.ErrObjectNotFound)
			//m.EXPECT().UpdateRoom(gomock.Any()).Return(nil)
			//act
			code := s.UpdateRoom([]byte(roomReq))

			// assert
			require.Equal(t, http.StatusNotFound, code)
		})

		t.Run("room with id doesn't exists", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_repository.NewMockDatabaseRepo(ctrl)
			s := NewRepo(m)

			m.EXPECT().GetRoomByID(id).Return(&room, nil)
			m.EXPECT().UpdateRoom(gomock.Any()).Return(repository.ErrInternalServer)
			//act
			code := s.UpdateRoom([]byte(roomReq))

			// assert
			require.Equal(t, http.StatusInternalServerError, code)
		})
	})

}

func Test_DeleteRoomWithAllReservations(t *testing.T) {
	var id int64 = 1

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := mock_repository.NewMockDatabaseRepo(ctrl)
		s := NewRepo(m)

		m.EXPECT().DeleteRoomByID(id).Return(nil)
		m.EXPECT().DeleteReservationsByRoomID(id).Return(nil)

		//act
		code := s.DeleteRoomWithAllReservations(id)

		// assert
		require.Equal(t, http.StatusOK, code)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("not found", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_repository.NewMockDatabaseRepo(ctrl)
			s := NewRepo(m)

			m.EXPECT().DeleteRoomByID(id).Return(repository.ErrObjectNotFound)

			//act
			code := s.DeleteRoomWithAllReservations(id)

			// assert
			require.Equal(t, http.StatusNotFound, code)
		})

		t.Run("internal error when delete room", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_repository.NewMockDatabaseRepo(ctrl)
			s := NewRepo(m)

			m.EXPECT().DeleteRoomByID(id).Return(repository.ErrInternalServer)

			//act
			code := s.DeleteRoomWithAllReservations(id)

			// assert
			require.Equal(t, http.StatusInternalServerError, code)
		})

		t.Run("internal error when delete reservations", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_repository.NewMockDatabaseRepo(ctrl)
			s := NewRepo(m)

			m.EXPECT().DeleteRoomByID(id).Return(nil)
			m.EXPECT().DeleteReservationsByRoomID(id).Return(repository.ErrInternalServer)

			//act
			code := s.DeleteRoomWithAllReservations(id)

			// assert
			require.Equal(t, http.StatusInternalServerError, code)
		})
	})
}

func Test_GetReservation(t *testing.T) {
	var (
		id int64 = 1
	)
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := mock_repository.NewMockDatabaseRepo(ctrl)
		s := NewRepo(m)

		m.EXPECT().GetReservationByID(id).Return(fixtures.Reservation().Valid().P(), nil)
		//act
		result, code := s.GetReservation(id)

		// assert
		require.Equal(t, http.StatusOK, code)
		assert.Equal(t, "{\"id\":1,\"start_date\":\"2023-11-08T00:00:00Z\",\"end_date\":\"2023-11-17T00:00:00Z\",\"room_id\":1,\"created_at\":\"0001-01-01T00:00:00Z\",\"updated_at\":\"0001-01-01T00:00:00Z\"}", string(result))
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := mock_repository.NewMockDatabaseRepo(ctrl)
		s := NewRepo(m)

		m.EXPECT().GetReservationByID(id).Return(nil, sql.ErrNoRows)
		//act
		result, code := s.GetReservation(id)

		// assert
		require.Equal(t, http.StatusInternalServerError, code)
		assert.Nil(t, result)
	})
}

func Test_CreateReservation(t *testing.T) {
	var (
		id     int64 = 1
		resReq       = "{\"start_date\":\"2023-09-15\", \"end_date\":\"2023-11-15\", \"room_id\":1}"
		room         = fixtures.Room().Valid().V()
	)
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := mock_repository.NewMockDatabaseRepo(ctrl)
		s := NewRepo(m)

		m.EXPECT().GetRoomByID(id).Return(&room, nil)
		m.EXPECT().InsertReservation(gomock.Any()).Return(id, nil)
		//act
		code := s.CreateReservation([]byte(resReq))

		// assert
		require.Equal(t, http.StatusOK, code)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()

		t.Run("room with id doesn't exists", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_repository.NewMockDatabaseRepo(ctrl)
			s := NewRepo(m)

			m.EXPECT().GetRoomByID(id).Return(nil, repository.ErrObjectNotFound)
			//act
			code := s.CreateReservation([]byte(resReq))

			// assert
			require.Equal(t, http.StatusNotFound, code)
		})

		t.Run("internal error when insert reservation", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_repository.NewMockDatabaseRepo(ctrl)
			s := NewRepo(m)

			m.EXPECT().GetRoomByID(id).Return(&room, nil)
			m.EXPECT().InsertReservation(gomock.Any()).Return(int64(0), repository.ErrInternalServer)
			//act
			code := s.CreateReservation([]byte(resReq))

			// assert
			require.Equal(t, http.StatusInternalServerError, code)
		})
	})
}

func Test_DeleteReservation(t *testing.T) {
	var id int64 = 1
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := mock_repository.NewMockDatabaseRepo(ctrl)
		s := NewRepo(m)

		m.EXPECT().DeleteReservationByID(id).Return(nil)
		//act
		code := s.DeleteReservation(id)

		// assert
		require.Equal(t, http.StatusOK, code)
	})
	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("not found", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_repository.NewMockDatabaseRepo(ctrl)
			s := NewRepo(m)

			m.EXPECT().DeleteReservationByID(id).Return(repository.ErrObjectNotFound)
			//act
			code := s.DeleteReservation(id)

			// assert
			require.Equal(t, http.StatusNotFound, code)
		})
		t.Run("internal server error", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_repository.NewMockDatabaseRepo(ctrl)
			s := NewRepo(m)

			m.EXPECT().DeleteReservationByID(id).Return(repository.ErrInternalServer)
			//act
			code := s.DeleteReservation(id)

			// assert
			require.Equal(t, http.StatusInternalServerError, code)
		})
	})

}

func Test_UpdateReservation(t *testing.T) {
	var (
		id     int64 = 1
		resReq       = "{\"id\":1, \"start_date\":\"2023-11-08\", \"end_date\":\"2023-11-17\", \"room_id\":1}"
		res          = fixtures.Reservation().Valid().P()
	)
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := mock_repository.NewMockDatabaseRepo(ctrl)
		s := NewRepo(m)

		m.EXPECT().GetReservationByID(id).Return(res, nil)
		m.EXPECT().UpdateReservation(gomock.Any()).Return(nil)
		//act
		code := s.UpdateReservation([]byte(resReq))

		// assert
		require.Equal(t, http.StatusOK, code)
	})

}
