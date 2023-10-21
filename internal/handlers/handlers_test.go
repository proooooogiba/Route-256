package handlers

import (
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"homework-3/internal/pkg/models"
	"homework-3/internal/pkg/repository"
	mock_repository "homework-3/internal/pkg/repository/mocks"
	"homework-3/tests/fixtures"
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
			getRoom, getReservations, err := s.GetRoomWithAllReservations(id)

			// assert
			require.Nil(t, err)
			assert.Equal(t, room, getRoom)
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
			getRoom, getReservations, err := s.GetRoomWithAllReservations(id)

			// assert
			require.Nil(t, err)
			assert.Equal(t, room, getRoom)
			assert.Equal(t, reservations, getReservations)
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
			getRoom, getReservations, err := s.GetRoomWithAllReservations(id)

			// assert
			require.ErrorIs(t, ErrRoomNotFound, err)
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
			getRoom, getReservations, err := s.GetRoomWithAllReservations(id)

			// assert
			require.ErrorIs(t, ErrInternalServer, err)
			assert.Nil(t, getReservations, getRoom)
		})
	})
}

func Test_CreateRoom(t *testing.T) {
	var (
		id   int64 = 1
		name       = "Lux"
		//roomReq       = "{\"name\":\"Lux\", \"cost\":1000.0}"
		room = fixtures.Room().Valid().V()
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
		err := s.CreateRoom(room)

		// assert
		require.Nil(t, err)
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
			err := s.CreateRoom(room)

			// assert
			require.ErrorIs(t, ErrRoomAlreadyExists, err)
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
			err := s.CreateRoom(room)

			// assert
			require.ErrorIs(t, ErrInternalServer, err)
		})
	})
}

func Test_UpdateRoom(t *testing.T) {
	var (
		id int64 = 1
		//name          = "Lux"
		//roomReq = "{\"id\":1, \"name\":\"Lux\", \"cost\":1000.0}"
		room = models.Room{
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
		err := s.UpdateRoom(room)

		// assert
		require.Nil(t, err)
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
			err := s.UpdateRoom(room)

			// assert
			require.ErrorIs(t, ErrRoomNotFound, err)
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
			err := s.UpdateRoom(room)

			// assert
			require.Equal(t, ErrInternalServer, err)
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
		err := s.DeleteRoomWithAllReservations(id)

		// assert
		require.Nil(t, err)
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
			err := s.DeleteRoomWithAllReservations(id)

			// assert
			require.ErrorIs(t, ErrRoomNotFound, err)
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
			err := s.DeleteRoomWithAllReservations(id)

			// assert
			require.ErrorIs(t, ErrInternalServer, err)
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
			err := s.DeleteRoomWithAllReservations(id)

			// assert
			require.ErrorIs(t, ErrInternalServer, err)
		})
	})
}

func Test_GetReservation(t *testing.T) {
	var (
		id  int64 = 1
		res       = fixtures.Reservation().Valid().P()
	)
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := mock_repository.NewMockDatabaseRepo(ctrl)
		s := NewRepo(m)

		m.EXPECT().GetReservationByID(id).Return(res, nil)
		//act
		reservation, err := s.GetReservation(id)

		// assert
		require.Nil(t, err)
		require.Equal(t, res, reservation)

		//assert.Equal(t, "{\"id\":1,\"start_date\":\"2023-11-08T00:00:00Z\",\"end_date\":\"2023-11-17T00:00:00Z\",\"room_id\":1,\"created_at\":\"0001-01-01T00:00:00Z\",\"updated_at\":\"0001-01-01T00:00:00Z\"}", string(result))
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := mock_repository.NewMockDatabaseRepo(ctrl)
		s := NewRepo(m)

		m.EXPECT().GetReservationByID(id).Return(nil, repository.ErrObjectNotFound)
		//act
		result, err := s.GetReservation(id)

		// assert
		require.ErrorIs(t, ErrReservationNotFound, err)
		assert.Nil(t, result)
	})
}

func Test_CreateReservation(t *testing.T) {
	var (
		id int64 = 1
		//resReq       = "{\"start_date\":\"2023-09-15\", \"end_date\":\"2023-11-15\", \"room_id\":1}"
		room = fixtures.Room().Valid().V()
		res  = fixtures.Reservation().Valid().V()
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
		err := s.CreateReservation(res)

		// assert
		require.Nil(t, err)
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
			err := s.CreateReservation(res)

			// assert
			require.ErrorIs(t, ErrRoomNotFound, err)
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
			err := s.CreateReservation(res)

			// assert
			require.ErrorIs(t, ErrInternalServer, err)
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
		err := s.DeleteReservation(id)

		// assert
		require.Nil(t, err)
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
			err := s.DeleteReservation(id)

			// assert
			require.ErrorIs(t, ErrReservationNotFound, err)
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
			err := s.DeleteReservation(id)

			// assert
			require.Equal(t, ErrInternalServer, err)
		})
	})
}

func Test_UpdateReservation(t *testing.T) {
	var (
		id int64 = 1
		//resReq       = "{\"id\":1, \"start_date\":\"2023-11-08\", \"end_date\":\"2023-11-17\", \"room_id\":1}"
		res = fixtures.Reservation().Valid().V()
	)
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := mock_repository.NewMockDatabaseRepo(ctrl)
		s := NewRepo(m)

		m.EXPECT().GetReservationByID(id).Return(&res, nil)
		m.EXPECT().UpdateReservation(gomock.Any()).Return(nil)
		//act
		err := s.UpdateReservation(res)

		// assert
		require.Nil(t, err)
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

			m.EXPECT().GetReservationByID(id).Return(nil, repository.ErrObjectNotFound)
			//m.EXPECT().UpdateReservation(gomock.Any()).Return(nil)
			//act
			err := s.UpdateReservation(res)

			// assert
			require.ErrorIs(t, ErrReservationNotFound, err)
		})
		t.Run("internal server error", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_repository.NewMockDatabaseRepo(ctrl)
			s := NewRepo(m)

			m.EXPECT().GetReservationByID(id).Return(&res, nil)
			m.EXPECT().UpdateReservation(gomock.Any()).Return(repository.ErrInternalServer)
			//act
			err := s.UpdateReservation(res)

			// assert
			require.ErrorIs(t, ErrInternalServer, err)
		})
	})
}
