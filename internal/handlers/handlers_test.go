package handlers

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"homework-3/internal/pkg/domain/mocks"
	"homework-3/internal/pkg/models"
	"homework-3/internal/pkg/parser"
	mock_parser "homework-3/internal/pkg/parser/mocks"
	"homework-3/internal/pkg/repository"
	"homework-3/internal/pkg/sender"
	mock_sender "homework-3/internal/pkg/sender/mocks"
	"homework-3/tests/fixtures"
	"net/http"
	"testing"
)

func TestService_CreateReservation(t *testing.T) {
	var (
		body                       = []byte(`{"start_date":"2023-11-08", "end_date":"2023-11-17", "room_id":1}`)
		invalidMarshalBody         = []byte(`{start_date":"2023-11-08", "eate":"2023-11-17", "room_id":1}`)
		invalidParseDateBody       = []byte(`{"start_date":"2023-101-088", "end_date":"2023-11-17", "room_id":1}`)
		res                        = fixtures.Reservation().CreateValid().V()
		sync                       = true
		resID                int64 = 1
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mock_repo.NewMockRepository(ctrl)
		mockSender := mock_sender.NewMockSender(ctrl)
		mockParser := mock_parser.NewMockParser(ctrl)
		s := NewService(mockRepo, mockSender, mockParser)

		mockParser.EXPECT().UnmarshalCreateReservationRequest(body).Return(res, nil)
		mockSender.EXPECT().Send(http.MethodPost, body, sync).Return(nil)
		mockRepo.EXPECT().CreateReservation(res).Return(resID, nil)

		//act
		_, err := s.CreateReservation(http.MethodPost, body, sync)

		// assert
		require.Nil(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("unmarshall", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mock_repo.NewMockRepository(ctrl)
			mockSender := mock_sender.NewMockSender(ctrl)
			mockParser := mock_parser.NewMockParser(ctrl)
			s := NewService(mockRepo, mockSender, mockParser)

			mockParser.EXPECT().UnmarshalCreateReservationRequest(invalidMarshalBody).Return(models.Reservation{}, parser.ErrUnmarshal)

			//act
			_, err := s.CreateReservation(http.MethodPost, invalidMarshalBody, sync)

			// assert
			require.ErrorIs(t, parser.ErrUnmarshal, err)
		})
		t.Run("parse date", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mock_repo.NewMockRepository(ctrl)
			mockSender := mock_sender.NewMockSender(ctrl)
			mockParser := mock_parser.NewMockParser(ctrl)
			s := NewService(mockRepo, mockSender, mockParser)

			mockParser.EXPECT().UnmarshalCreateReservationRequest(invalidParseDateBody).Return(models.Reservation{}, parser.ErrParseDate)

			//act
			_, err := s.CreateReservation(http.MethodPost, invalidParseDateBody, sync)

			// assert
			require.ErrorIs(t, parser.ErrParseDate, err)
		})
		t.Run("db operation", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mock_repo.NewMockRepository(ctrl)
			mockSender := mock_sender.NewMockSender(ctrl)
			mockParser := mock_parser.NewMockParser(ctrl)
			s := NewService(mockRepo, mockSender, mockParser)

			mockParser.EXPECT().UnmarshalCreateReservationRequest(body).Return(res, nil)
			mockSender.EXPECT().Send(http.MethodPost, body, sync).Return(nil)
			mockRepo.EXPECT().CreateReservation(res).Return(int64(0), repository.ErrInternalServer)
			//act
			_, err := s.CreateReservation(http.MethodPost, body, sync)

			// assert
			require.ErrorIs(t, repository.ErrInternalServer, err)
		})
		t.Run("kafka_test sender", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mock_repo.NewMockRepository(ctrl)
			mockSender := mock_sender.NewMockSender(ctrl)
			mockParser := mock_parser.NewMockParser(ctrl)
			s := NewService(mockRepo, mockSender, mockParser)

			mockParser.EXPECT().UnmarshalCreateReservationRequest(body).Return(res, nil)
			mockSender.EXPECT().Send(http.MethodPost, body, sync).Return(sender.ErrSendSyncMessage)

			//act
			_, err := s.CreateReservation(http.MethodPost, body, sync)

			// assert
			require.ErrorIs(t, sender.ErrSendSyncMessage, err)
		})
	})
}

func TestService_GetReservation(t *testing.T) {
	var (
		id   int64 = 1
		body       = []byte("")
		res        = fixtures.Reservation().CreateValid().P()
		sync       = true
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mock_repo.NewMockRepository(ctrl)
		mockSender := mock_sender.NewMockSender(ctrl)
		mockParser := mock_parser.NewMockParser(ctrl)
		s := NewService(mockRepo, mockSender, mockParser)

		mockSender.EXPECT().Send(http.MethodGet, body, sync).Return(nil)
		mockRepo.EXPECT().GetReservation(id).Return(res, nil)

		//act
		getRes, err := s.GetReservation(http.MethodGet, id, sync)

		// assert
		require.Nil(t, err)
		require.Equal(t, res, getRes)

	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("db operation", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mock_repo.NewMockRepository(ctrl)
			mockSender := mock_sender.NewMockSender(ctrl)
			mockParser := mock_parser.NewMockParser(ctrl)
			s := NewService(mockRepo, mockSender, mockParser)

			mockSender.EXPECT().Send(http.MethodGet, body, sync).Return(nil)
			mockRepo.EXPECT().GetReservation(id).Return(nil, repository.ErrObjectNotFound)

			//act
			getRes, err := s.GetReservation(http.MethodGet, id, sync)

			// assert
			require.ErrorIs(t, err, repository.ErrObjectNotFound)
			require.Zero(t, getRes)
		})
		t.Run("kafka_test sender", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mock_repo.NewMockRepository(ctrl)
			mockSender := mock_sender.NewMockSender(ctrl)
			mockParser := mock_parser.NewMockParser(ctrl)
			s := NewService(mockRepo, mockSender, mockParser)

			mockSender.EXPECT().Send(http.MethodGet, body, sync).Return(sender.ErrSendSyncMessage)

			//act
			getRes, err := s.GetReservation(http.MethodGet, id, sync)

			// assert
			require.ErrorIs(t, err, sender.ErrSendSyncMessage)
			require.Zero(t, getRes)
		})
	})
}

func TestService_UpdateReservation(t *testing.T) {
	var (
		body                 = []byte(`{"id":1, "start_date":"2023-11-08", "end_date":"2023-11-17", "room_id":1}`)
		invalidMarshalBody   = []byte(`{"":1, start_date":"2023-11-08", "eate":"2023-11-17", "room":1}`)
		invalidParseDateBody = []byte(`{"id":1, "start_date":"2023-101-088", "end_date":"2023-11-17", "room_id":1}`)
		res                  = fixtures.Reservation().CreateValid().V()
		sync                 = true
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mock_repo.NewMockRepository(ctrl)
		mockSender := mock_sender.NewMockSender(ctrl)
		mockParser := mock_parser.NewMockParser(ctrl)
		s := NewService(mockRepo, mockSender, mockParser)

		mockParser.EXPECT().UnmarshalUpdateReservationRequest(body).Return(res, nil)
		mockSender.EXPECT().Send(http.MethodPut, body, sync).Return(nil)
		mockRepo.EXPECT().UpdateReservation(res).Return(nil)

		//act
		err := s.UpdateReservation(http.MethodPut, body, sync)

		// assert
		require.Nil(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("unmarshall", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mock_repo.NewMockRepository(ctrl)
			mockSender := mock_sender.NewMockSender(ctrl)
			mockParser := mock_parser.NewMockParser(ctrl)
			s := NewService(mockRepo, mockSender, mockParser)

			mockParser.EXPECT().UnmarshalUpdateReservationRequest(invalidMarshalBody).Return(models.Reservation{}, parser.ErrUnmarshal)

			//act
			err := s.UpdateReservation(http.MethodPut, invalidMarshalBody, sync)

			// assert
			require.ErrorIs(t, parser.ErrUnmarshal, err)
		})
		t.Run("parse date", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mock_repo.NewMockRepository(ctrl)
			mockSender := mock_sender.NewMockSender(ctrl)
			mockParser := mock_parser.NewMockParser(ctrl)
			s := NewService(mockRepo, mockSender, mockParser)

			mockParser.EXPECT().UnmarshalUpdateReservationRequest(invalidParseDateBody).Return(models.Reservation{}, parser.ErrParseDate)

			//act
			err := s.UpdateReservation(http.MethodPut, invalidParseDateBody, sync)

			// assert
			require.ErrorIs(t, parser.ErrParseDate, err)
		})
		t.Run("db operation", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mock_repo.NewMockRepository(ctrl)
			mockSender := mock_sender.NewMockSender(ctrl)
			mockParser := mock_parser.NewMockParser(ctrl)
			s := NewService(mockRepo, mockSender, mockParser)

			mockParser.EXPECT().UnmarshalUpdateReservationRequest(body).Return(res, nil)
			mockSender.EXPECT().Send(http.MethodPut, body, sync).Return(nil)
			mockRepo.EXPECT().UpdateReservation(res).Return(repository.ErrInternalServer)

			//act
			err := s.UpdateReservation(http.MethodPut, body, sync)

			// assert
			require.ErrorIs(t, repository.ErrInternalServer, err)
		})
		t.Run("kafka_test sender", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mock_repo.NewMockRepository(ctrl)
			mockSender := mock_sender.NewMockSender(ctrl)
			mockParser := mock_parser.NewMockParser(ctrl)
			s := NewService(mockRepo, mockSender, mockParser)

			mockParser.EXPECT().UnmarshalUpdateReservationRequest(body).Return(res, nil)
			mockSender.EXPECT().Send(http.MethodPut, body, sync).Return(sender.ErrSendSyncMessage)

			//act
			err := s.UpdateReservation(http.MethodPut, body, sync)

			// assert
			require.ErrorIs(t, sender.ErrSendSyncMessage, err)
		})
	})
}

func TestService_DeleteReservation(t *testing.T) {
	var (
		id   int64 = 1
		body       = []byte("")
		sync       = true
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mock_repo.NewMockRepository(ctrl)
		mockSender := mock_sender.NewMockSender(ctrl)
		mockParser := mock_parser.NewMockParser(ctrl)
		s := NewService(mockRepo, mockSender, mockParser)

		mockSender.EXPECT().Send(http.MethodDelete, body, sync).Return(nil)
		mockRepo.EXPECT().DeleteReservation(id).Return(nil)

		//act
		err := s.DeleteReservation(http.MethodDelete, id, sync)

		// assert
		require.Nil(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("db operation", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mock_repo.NewMockRepository(ctrl)
			mockSender := mock_sender.NewMockSender(ctrl)
			mockParser := mock_parser.NewMockParser(ctrl)
			s := NewService(mockRepo, mockSender, mockParser)

			mockSender.EXPECT().Send(http.MethodDelete, body, sync).Return(nil)
			mockRepo.EXPECT().DeleteReservation(id).Return(repository.ErrInternalServer)

			//act
			err := s.DeleteReservation(http.MethodDelete, id, sync)

			// assert
			require.ErrorIs(t, err, repository.ErrInternalServer)
		})
		t.Run("kafka_test sender", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mock_repo.NewMockRepository(ctrl)
			mockSender := mock_sender.NewMockSender(ctrl)
			mockParser := mock_parser.NewMockParser(ctrl)
			s := NewService(mockRepo, mockSender, mockParser)

			mockSender.EXPECT().Send(http.MethodDelete, body, sync).Return(sender.ErrSendSyncMessage)

			//act
			err := s.DeleteReservation(http.MethodDelete, id, sync)

			// assert
			require.ErrorIs(t, err, sender.ErrSendSyncMessage)
		})
	})
}

func TestService_CreateRoom(t *testing.T) {
	var (
		body                     = []byte(`{"name":"Lux", "cost":1000.0}`)
		invalidMarshalBody       = []byte(`{"na":"Lux", "cost":100a0.0}`)
		room                     = fixtures.Room().CreateValid().V()
		sync                     = true
		roomID             int64 = 1
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mock_repo.NewMockRepository(ctrl)
		mockSender := mock_sender.NewMockSender(ctrl)
		mockParser := mock_parser.NewMockParser(ctrl)
		s := NewService(mockRepo, mockSender, mockParser)

		mockParser.EXPECT().UnmarshalCreateRoomRequest(body).Return(room, nil)
		mockSender.EXPECT().Send(http.MethodPost, body, sync).Return(nil)
		mockRepo.EXPECT().CreateRoom(room).Return(roomID, nil)

		//act
		_, err := s.CreateRoom(http.MethodPost, body, sync)

		// assert
		require.Nil(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("unmarshall", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mock_repo.NewMockRepository(ctrl)
			mockSender := mock_sender.NewMockSender(ctrl)
			mockParser := mock_parser.NewMockParser(ctrl)
			s := NewService(mockRepo, mockSender, mockParser)

			mockParser.EXPECT().UnmarshalCreateRoomRequest(invalidMarshalBody).Return(models.Room{}, parser.ErrUnmarshal)

			//act
			_, err := s.CreateRoom(http.MethodPost, invalidMarshalBody, sync)

			// assert
			require.ErrorIs(t, parser.ErrUnmarshal, err)
		})
		t.Run("db operation", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mock_repo.NewMockRepository(ctrl)
			mockSender := mock_sender.NewMockSender(ctrl)
			mockParser := mock_parser.NewMockParser(ctrl)
			s := NewService(mockRepo, mockSender, mockParser)

			mockParser.EXPECT().UnmarshalCreateRoomRequest(body).Return(room, nil)
			mockSender.EXPECT().Send(http.MethodPost, body, sync).Return(nil)
			mockRepo.EXPECT().CreateRoom(room).Return(int64(0), repository.ErrInternalServer)

			//act
			_, err := s.CreateRoom(http.MethodPost, body, sync)

			// assert
			require.ErrorIs(t, repository.ErrInternalServer, err)
		})
		t.Run("kafka_test sender", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mock_repo.NewMockRepository(ctrl)
			mockSender := mock_sender.NewMockSender(ctrl)
			mockParser := mock_parser.NewMockParser(ctrl)
			s := NewService(mockRepo, mockSender, mockParser)

			mockParser.EXPECT().UnmarshalCreateRoomRequest(body).Return(room, nil)
			mockSender.EXPECT().Send(http.MethodPost, body, sync).Return(sender.ErrSendSyncMessage)

			//act
			_, err := s.CreateRoom(http.MethodPost, body, sync)

			// assert
			require.ErrorIs(t, sender.ErrSendSyncMessage, err)
		})
	})
}

func TestService_GetRoomWithAllReservations(t *testing.T) {
	var (
		roomID int64 = 1
		body         = []byte("")
		sync         = true
		room         = fixtures.Room().Valid().P()
	)

	reservations := make([]*models.Reservation, 2, 4)
	reservations = append(reservations, fixtures.Reservation().Valid().P())
	reservations = append(reservations, fixtures.Reservation().Valid2().P())

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mock_repo.NewMockRepository(ctrl)
		mockSender := mock_sender.NewMockSender(ctrl)
		mockParser := mock_parser.NewMockParser(ctrl)
		s := NewService(mockRepo, mockSender, mockParser)

		mockSender.EXPECT().Send(http.MethodGet, body, sync).Return(nil)
		mockRepo.EXPECT().GetRoomWithAllReservations(roomID).Return(room, reservations, nil)

		//act
		getRoom, getReservations, err := s.GetRoomWithAllReservations(http.MethodGet, roomID, sync)

		// assert
		require.Nil(t, err)
		require.Equal(t, getRoom, room)
		require.Equal(t, getReservations, reservations)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("db operation", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mock_repo.NewMockRepository(ctrl)
			mockSender := mock_sender.NewMockSender(ctrl)
			mockParser := mock_parser.NewMockParser(ctrl)
			s := NewService(mockRepo, mockSender, mockParser)

			mockSender.EXPECT().Send(http.MethodGet, body, sync).Return(nil)
			mockRepo.EXPECT().GetRoomWithAllReservations(roomID).Return(room, reservations, repository.ErrInternalServer)

			//act
			getRoom, getReservations, err := s.GetRoomWithAllReservations(http.MethodGet, roomID, sync)

			// assert
			require.ErrorIs(t, err, repository.ErrInternalServer)
			require.Nil(t, getRoom, getReservations)
		})
		t.Run("kafka_test sender", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mock_repo.NewMockRepository(ctrl)
			mockSender := mock_sender.NewMockSender(ctrl)
			mockParser := mock_parser.NewMockParser(ctrl)
			s := NewService(mockRepo, mockSender, mockParser)

			mockSender.EXPECT().Send(http.MethodGet, body, sync).Return(sender.ErrSendSyncMessage)

			//act
			getRoom, getReservations, err := s.GetRoomWithAllReservations(http.MethodGet, roomID, sync)

			// assert
			require.ErrorIs(t, err, sender.ErrSendSyncMessage)
			require.Nil(t, getRoom, getReservations)
		})
	})
}

func TestService_UpdateRoom(t *testing.T) {
	var (
		body               = []byte(`{"id":1, "name":"Lux", "cost":1000.0}`)
		invalidMarshalBody = []byte(`{"na":"Lux", "cost":100a0.0}`)
		room               = fixtures.Room().Valid().V()
		sync               = true
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mock_repo.NewMockRepository(ctrl)
		mockSender := mock_sender.NewMockSender(ctrl)
		mockParser := mock_parser.NewMockParser(ctrl)
		s := NewService(mockRepo, mockSender, mockParser)

		mockParser.EXPECT().UnmarshalUpdateRoomRequest(body).Return(room, nil)
		mockSender.EXPECT().Send(http.MethodPut, body, sync).Return(nil)
		mockRepo.EXPECT().UpdateRoom(room).Return(nil)

		//act
		err := s.UpdateRoom(http.MethodPut, body, sync)

		// assert
		require.Nil(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("unmarshall", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mock_repo.NewMockRepository(ctrl)
			mockSender := mock_sender.NewMockSender(ctrl)
			mockParser := mock_parser.NewMockParser(ctrl)
			s := NewService(mockRepo, mockSender, mockParser)

			mockParser.EXPECT().UnmarshalUpdateRoomRequest(invalidMarshalBody).Return(models.Room{}, parser.ErrUnmarshal)

			//act
			err := s.UpdateRoom(http.MethodPut, invalidMarshalBody, sync)

			// assert
			require.ErrorIs(t, parser.ErrUnmarshal, err)
		})
		t.Run("db operation", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mock_repo.NewMockRepository(ctrl)
			mockSender := mock_sender.NewMockSender(ctrl)
			mockParser := mock_parser.NewMockParser(ctrl)
			s := NewService(mockRepo, mockSender, mockParser)

			mockParser.EXPECT().UnmarshalUpdateRoomRequest(body).Return(room, nil)
			mockSender.EXPECT().Send(http.MethodPut, body, sync).Return(nil)
			mockRepo.EXPECT().UpdateRoom(room).Return(repository.ErrInternalServer)

			//act
			err := s.UpdateRoom(http.MethodPut, body, sync)

			// assert
			require.ErrorIs(t, repository.ErrInternalServer, err)
		})
		t.Run("kafka_test sender", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mock_repo.NewMockRepository(ctrl)
			mockSender := mock_sender.NewMockSender(ctrl)
			mockParser := mock_parser.NewMockParser(ctrl)
			s := NewService(mockRepo, mockSender, mockParser)

			mockParser.EXPECT().UnmarshalUpdateRoomRequest(body).Return(room, nil)
			mockSender.EXPECT().Send(http.MethodPut, body, sync).Return(sender.ErrSendSyncMessage)

			//act
			err := s.UpdateRoom(http.MethodPut, body, sync)

			// assert
			require.ErrorIs(t, sender.ErrSendSyncMessage, err)
		})
	})
}

func TestService_DeleteRoomWithAllReservations(t *testing.T) {
	var (
		roomID int64 = 1
		body         = []byte("")
		sync         = true
	)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		// arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mock_repo.NewMockRepository(ctrl)
		mockSender := mock_sender.NewMockSender(ctrl)
		mockParser := mock_parser.NewMockParser(ctrl)
		s := NewService(mockRepo, mockSender, mockParser)

		mockSender.EXPECT().Send(http.MethodDelete, body, sync).Return(nil)
		mockRepo.EXPECT().DeleteRoomWithAllReservations(roomID).Return(nil)

		//act
		err := s.DeleteRoomWithAllReservations(http.MethodDelete, roomID, sync)

		// assert
		require.Nil(t, err)
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("db operation", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mock_repo.NewMockRepository(ctrl)
			mockSender := mock_sender.NewMockSender(ctrl)
			mockParser := mock_parser.NewMockParser(ctrl)
			s := NewService(mockRepo, mockSender, mockParser)

			mockSender.EXPECT().Send(http.MethodDelete, body, sync).Return(nil)
			mockRepo.EXPECT().DeleteRoomWithAllReservations(roomID).Return(repository.ErrInternalServer)

			//act
			err := s.DeleteRoomWithAllReservations(http.MethodDelete, roomID, sync)

			// assert
			require.ErrorIs(t, err, repository.ErrInternalServer)
		})
		t.Run("kafka_test sender", func(t *testing.T) {
			t.Parallel()

			// arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockRepo := mock_repo.NewMockRepository(ctrl)
			mockSender := mock_sender.NewMockSender(ctrl)
			mockParser := mock_parser.NewMockParser(ctrl)
			s := NewService(mockRepo, mockSender, mockParser)

			mockSender.EXPECT().Send(http.MethodDelete, body, sync).Return(sender.ErrSendSyncMessage)

			//act
			err := s.DeleteRoomWithAllReservations(http.MethodDelete, roomID, sync)

			// assert
			require.ErrorIs(t, err, sender.ErrSendSyncMessage)
		})
	})
}
