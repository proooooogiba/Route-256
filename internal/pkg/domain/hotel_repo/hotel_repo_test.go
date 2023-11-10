package hotel_repo

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"homework-3/internal/pkg/domain"
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
		ctx        = context.Background()
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

			m.EXPECT().GetRoomByID(gomock.Any(), id).Return(room, nil)
			m.EXPECT().GetReservationsByRoomID(gomock.Any(), id).Return(nil, nil)

			//act
			getRoom, getReservations, err := s.GetRoomWithAllReservations(ctx, id)

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

			m.EXPECT().GetRoomByID(gomock.Any(), id).Return(room, nil)
			m.EXPECT().GetReservationsByRoomID(gomock.Any(), id).Return(reservations, nil)
			//act
			getRoom, getReservations, err := s.GetRoomWithAllReservations(ctx, id)

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

			m.EXPECT().GetRoomByID(gomock.Any(), id).Return(nil, repository.ErrObjectNotFound)
			//act
			getRoom, getReservations, err := s.GetRoomWithAllReservations(ctx, id)

			// assert
			require.ErrorIs(t, domain.ErrRoomNotFound, err)
			assert.Nil(t, getRoom, getReservations)
		})

		t.Run("internal error", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_repository.NewMockDatabaseRepo(ctrl)
			s := NewRepo(m)

			m.EXPECT().GetRoomByID(gomock.Any(), id).Return(room, nil)
			m.EXPECT().GetReservationsByRoomID(gomock.Any(), id).Return(nil, errors.New("Error while getting reservations by room id"))
			//act
			getRoom, getReservations, err := s.GetRoomWithAllReservations(ctx, id)

			// assert
			require.ErrorIs(t, domain.ErrInternalServer, err)
			assert.Nil(t, getReservations, getRoom)
		})
	})
}

func Test_CreateRoom(t *testing.T) {
	var (
		id   int64 = 1
		name       = "Lux"
		room       = fixtures.Room().Valid().V()
		ctx        = context.Background()
	)
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := mock_repository.NewMockDatabaseRepo(ctrl)
		s := NewRepo(m)

		m.EXPECT().GetRoomByName(gomock.Any(), name).Return(nil, repository.ErrObjectNotFound)
		m.EXPECT().InsertRoom(gomock.Any(), &room).Return(id, nil)
		//act
		roomID, err := s.CreateRoom(ctx, room)

		// assert
		require.Nil(t, err)
		require.Equal(t, id, roomID)
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

			m.EXPECT().GetRoomByName(gomock.Any(), name).Return(nil, nil)
			//act
			roomID, err := s.CreateRoom(ctx, room)

			// assert
			require.ErrorIs(t, domain.ErrRoomAlreadyExists, err)
			require.Zero(t, roomID)
		})
		t.Run("room with following name already exists", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_repository.NewMockDatabaseRepo(ctrl)
			s := NewRepo(m)

			m.EXPECT().GetRoomByName(gomock.Any(), name).Return(nil, repository.ErrObjectNotFound)
			m.EXPECT().InsertRoom(gomock.Any(), &room).Return(id, repository.ErrInternalServer)

			//act
			roomID, err := s.CreateRoom(ctx, room)

			// assert
			require.ErrorIs(t, domain.ErrInternalServer, err)
			require.Zero(t, roomID)
		})
	})
}

func Test_UpdateRoom(t *testing.T) {
	var (
		id   int64 = 1
		room       = models.Room{
			ID:   1,
			Name: "Lux",
			Cost: 1000.0,
		}
		ctx = context.Background()
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := mock_repository.NewMockDatabaseRepo(ctrl)
		s := NewRepo(m)

		m.EXPECT().GetRoomByID(gomock.Any(), id).Return(&room, nil)
		m.EXPECT().UpdateRoom(gomock.Any(), &room).Return(nil)
		//act
		err := s.UpdateRoom(ctx, room)

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

			m.EXPECT().GetRoomByID(gomock.Any(), id).Return(nil, repository.ErrObjectNotFound)
			//act
			err := s.UpdateRoom(ctx, room)

			// assert
			require.ErrorIs(t, domain.ErrRoomNotFound, err)
		})

		t.Run("room with id doesn't exists", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_repository.NewMockDatabaseRepo(ctrl)
			s := NewRepo(m)

			m.EXPECT().GetRoomByID(gomock.Any(), id).Return(&room, nil)
			m.EXPECT().UpdateRoom(gomock.Any(), &room).Return(repository.ErrInternalServer)
			//act
			err := s.UpdateRoom(ctx, room)

			// assert
			require.Equal(t, domain.ErrInternalServer, err)
		})
	})

}

func Test_DeleteRoomWithAllReservations(t *testing.T) {
	var (
		id  int64 = 1
		ctx       = context.Background()
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := mock_repository.NewMockDatabaseRepo(ctrl)
		s := NewRepo(m)

		m.EXPECT().DeleteRoomByID(gomock.Any(), id).Return(nil)
		m.EXPECT().DeleteReservationsByRoomID(gomock.Any(), id).Return(nil)

		//act
		err := s.DeleteRoomWithAllReservations(ctx, id)

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

			m.EXPECT().DeleteRoomByID(gomock.Any(), id).Return(repository.ErrObjectNotFound)
			//act
			err := s.DeleteRoomWithAllReservations(ctx, id)

			// assert
			require.ErrorIs(t, domain.ErrRoomNotFound, err)
		})

		t.Run("internal error when delete room", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_repository.NewMockDatabaseRepo(ctrl)
			s := NewRepo(m)

			m.EXPECT().DeleteRoomByID(gomock.Any(), id).Return(repository.ErrInternalServer)
			//act
			err := s.DeleteRoomWithAllReservations(ctx, id)

			// assert
			require.ErrorIs(t, domain.ErrInternalServer, err)
		})

		t.Run("internal error when delete reservations", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_repository.NewMockDatabaseRepo(ctrl)
			s := NewRepo(m)

			m.EXPECT().DeleteRoomByID(gomock.Any(), id).Return(nil)
			m.EXPECT().DeleteReservationsByRoomID(gomock.Any(), id).Return(repository.ErrInternalServer)

			//act
			err := s.DeleteRoomWithAllReservations(ctx, id)

			// assert
			require.ErrorIs(t, domain.ErrInternalServer, err)
		})
	})
}

func Test_GetReservation(t *testing.T) {
	var (
		id  int64 = 1
		res       = fixtures.Reservation().Valid().P()
		ctx       = context.Background()
	)
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := mock_repository.NewMockDatabaseRepo(ctrl)
		s := NewRepo(m)

		m.EXPECT().GetReservationByID(gomock.Any(), id).Return(res, nil)
		//act
		reservation, err := s.GetReservation(ctx, id)

		// assert
		require.Nil(t, err)
		require.Equal(t, res, reservation)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := mock_repository.NewMockDatabaseRepo(ctrl)
		s := NewRepo(m)

		m.EXPECT().GetReservationByID(gomock.Any(), id).Return(nil, repository.ErrObjectNotFound)
		//act
		result, err := s.GetReservation(ctx, id)

		// assert
		require.ErrorIs(t, domain.ErrReservationNotFound, err)
		assert.Nil(t, result)
	})
}

func Test_CreateReservation(t *testing.T) {
	var (
		id   int64 = 1
		room       = fixtures.Room().Valid().V()
		res        = fixtures.Reservation().Valid().V()
		ctx        = context.Background()
	)
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := mock_repository.NewMockDatabaseRepo(ctrl)
		s := NewRepo(m)

		m.EXPECT().GetRoomByID(gomock.Any(), id).Return(&room, nil)
		m.EXPECT().InsertReservation(gomock.Any(), &res).Return(id, nil)
		//act
		_, err := s.CreateReservation(ctx, res)

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

			m.EXPECT().GetRoomByID(gomock.Any(), id).Return(nil, repository.ErrObjectNotFound)
			//act
			_, err := s.CreateReservation(ctx, res)

			// assert
			require.ErrorIs(t, domain.ErrRoomNotFound, err)
		})

		t.Run("internal error when insert reservation", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_repository.NewMockDatabaseRepo(ctrl)
			s := NewRepo(m)

			m.EXPECT().GetRoomByID(gomock.Any(), id).Return(&room, nil)
			m.EXPECT().InsertReservation(gomock.Any(), &res).Return(int64(0), repository.ErrInternalServer)
			//act
			_, err := s.CreateReservation(ctx, res)

			// assert
			require.ErrorIs(t, domain.ErrInternalServer, err)
		})
	})
}

func Test_DeleteReservation(t *testing.T) {
	var (
		id  int64 = 1
		ctx       = context.Background()
	)
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := mock_repository.NewMockDatabaseRepo(ctrl)
		s := NewRepo(m)

		m.EXPECT().DeleteReservationByID(gomock.Any(), id).Return(nil)
		//act
		err := s.DeleteReservation(ctx, id)

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

			m.EXPECT().DeleteReservationByID(gomock.Any(), id).Return(repository.ErrObjectNotFound)
			//act
			err := s.DeleteReservation(ctx, id)

			// assert
			require.ErrorIs(t, domain.ErrReservationNotFound, err)
		})
		t.Run("internal server error", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_repository.NewMockDatabaseRepo(ctrl)
			s := NewRepo(m)

			m.EXPECT().DeleteReservationByID(gomock.Any(), id).Return(repository.ErrInternalServer)
			//act
			err := s.DeleteReservation(ctx, id)

			// assert
			require.Equal(t, domain.ErrInternalServer, err)
		})
	})
}

func Test_UpdateReservation(t *testing.T) {
	var (
		id  int64 = 1
		res       = fixtures.Reservation().Valid().V()
		ctx       = context.Background()
	)
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		m := mock_repository.NewMockDatabaseRepo(ctrl)
		s := NewRepo(m)

		m.EXPECT().GetReservationByID(gomock.Any(), id).Return(&res, nil)
		m.EXPECT().UpdateReservation(gomock.Any(), &res).Return(nil)
		//act
		err := s.UpdateReservation(ctx, res)

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

			m.EXPECT().GetReservationByID(gomock.Any(), id).Return(nil, repository.ErrObjectNotFound)
			//act
			err := s.UpdateReservation(ctx, res)

			// assert
			require.ErrorIs(t, domain.ErrReservationNotFound, err)
		})
		t.Run("internal server error", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			m := mock_repository.NewMockDatabaseRepo(ctrl)
			s := NewRepo(m)

			m.EXPECT().GetReservationByID(gomock.Any(), id).Return(&res, nil)
			m.EXPECT().UpdateReservation(gomock.Any(), &res).Return(repository.ErrInternalServer)
			//act
			err := s.UpdateReservation(ctx, res)

			// assert
			require.ErrorIs(t, domain.ErrInternalServer, err)
		})
	})
}
