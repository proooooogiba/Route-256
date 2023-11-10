//go:build integration
// +build integration

package tests

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"homework-3/internal/handlers"
	"homework-3/internal/pkg/domain/hotel_repo"
	"homework-3/internal/pkg/parser/parser_request"
	"homework-3/internal/pkg/repository/dbrepo"
	"homework-3/internal/producer"
	"io"
	"net/http"
	"os"
	"testing"
)

func Test_CreateRoom(t *testing.T) {
	t.Parallel()
	var (
		createBody = []byte(`{"name":"Lux", "cost":1000.0}`)
		topic      = viper.GetString("TOPIC")
		outBody    = "\"{\\\"name\\\":\\\"Lux\\\", \\\"cost\\\":1000.0}\""
		sync       = true
	)
	t.Run("success", func(t *testing.T) {
		db.SetUp(t)
		defer db.TearDown()

		tProducer.SetUp(t)
		defer tProducer.TearDown()

		tConsumer.SetUp(t)
		defer tConsumer.TearDown()

		wantMessages := 1
		messagesChan := make(chan bool, wantMessages)
		tConsumer.Consumer.StartConsume(topic, wantMessages, messagesChan)

		hotelService := handlers.NewService(
			hotel_repo.NewRepo(dbrepo.NewPostgresRepo(db.DB)),
			producer.NewKafkaSender(tProducer.Producer, topic),
			parser_request.NewRequestParser(),
		)

		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		//act
		_, err := hotelService.CreateRoom(http.MethodPost, createBody, sync)
		require.Nil(t, err)

		for i := 0; i < wantMessages; i++ {
			<-messagesChan
		}
		w.Close()
		out, _ := io.ReadAll(r)
		os.Stdout = rescueStdout

		//assert
		require.Contains(t, string(out), `"method":"POST"`)
		require.Contains(t, string(out), outBody)
	})
}

func Test_GetRoomWithAllReservations(t *testing.T) {
	t.Parallel()
	var (
		createBody    = []byte("{\"name\":\"Lux\", \"cost\":1000.0}")
		topic         = viper.GetString("TOPIC")
		createOutBody = "\"{\\\"name\\\":\\\"Lux\\\", \\\"cost\\\":1000.0}\""
		getOutBody    = "\"body\":\"\""
		sync          = true
		wantMessages  = 2
	)

	t.Run("success", func(t *testing.T) {
		db.SetUp(t)
		defer db.TearDown()

		tProducer.SetUp(t)
		defer tProducer.TearDown()

		tConsumer.SetUp(t)
		defer tConsumer.TearDown()

		messagesChan := make(chan bool, wantMessages)
		tConsumer.Consumer.StartConsume(topic, wantMessages, messagesChan)

		hotelService := handlers.NewService(
			hotel_repo.NewRepo(dbrepo.NewPostgresRepo(db.DB)),
			producer.NewKafkaSender(tProducer.Producer, topic),
			parser_request.NewRequestParser(),
		)

		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		//act
		roomID, err := hotelService.CreateRoom(http.MethodPost, createBody, sync)
		require.Nil(t, err)

		_, _, err = hotelService.GetRoomWithAllReservations(http.MethodGet, roomID, sync)
		require.Nil(t, err)

		for i := 0; i < wantMessages; i++ {
			<-messagesChan
		}
		w.Close()
		out, _ := io.ReadAll(r)
		os.Stdout = rescueStdout

		//assert
		require.Contains(t, string(out), `"method":"GET"`)
		require.Contains(t, string(out), createOutBody)
		require.Contains(t, string(out), getOutBody)
	})
}

func Test_UpdateRoomKafka(t *testing.T) {
	t.Parallel()
	var (
		createBody       = []byte("{\"name\":\"Lux\", \"cost\":1000.0}")
		topic            = viper.GetString("TOPIC")
		updateStrBody    = "{\"id\":%v, \"name\":\"Lux\", \"cost\":1000.0}"
		createOutBody    = "\"{\\\"name\\\":\\\"Lux\\\", \\\"cost\\\":1000.0}\""
		updateOutStrBody = "\"{\\\"id\\\":%v, \\\"name\\\":\\\"Lux\\\", \\\"cost\\\":1000.0}\""
		sync             = true
		wantMessages     = 2
	)

	t.Run("success", func(t *testing.T) {
		db.SetUp(t)
		defer db.TearDown()

		tProducer.SetUp(t)
		defer tProducer.TearDown()

		tConsumer.SetUp(t)
		defer tConsumer.TearDown()

		messagesChan := make(chan bool, wantMessages)
		tConsumer.Consumer.StartConsume(topic, wantMessages, messagesChan)

		hotelService := handlers.NewService(
			hotel_repo.NewRepo(dbrepo.NewPostgresRepo(db.DB)),
			producer.NewKafkaSender(tProducer.Producer, topic),
			parser_request.NewRequestParser(),
		)

		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		//act
		roomID, err := hotelService.CreateRoom(http.MethodPost, createBody, sync)
		require.Nil(t, err)

		updateBody := []byte(fmt.Sprintf(updateStrBody, roomID))
		err = hotelService.UpdateRoom(http.MethodPut, updateBody, sync)
		require.Nil(t, err)

		for i := 0; i < wantMessages; i++ {
			<-messagesChan
		}
		w.Close()
		out, _ := io.ReadAll(r)
		os.Stdout = rescueStdout

		//assert
		require.Contains(t, string(out), `"method":"PUT"`)
		require.Contains(t, string(out), createOutBody)
		require.Contains(t, string(out), fmt.Sprintf(updateOutStrBody, roomID))
	})
}

func Test_DeleteRoomWithAllReservations(t *testing.T) {
	t.Parallel()
	var (
		createBody    = []byte("{\"name\":\"Lux\", \"cost\":1000.0}")
		topic         = viper.GetString("TOPIC")
		createOutBody = "\"{\\\"name\\\":\\\"Lux\\\", \\\"cost\\\":1000.0}\""
		deleteOutBody = "\"body\":\"\""
		sync          = true
		wantMessages  = 2
	)

	t.Run("success", func(t *testing.T) {
		db.SetUp(t)
		defer db.TearDown()

		tProducer.SetUp(t)
		defer tProducer.TearDown()

		tConsumer.SetUp(t)
		defer tConsumer.TearDown()

		messagesChan := make(chan bool, wantMessages)
		tConsumer.Consumer.StartConsume(topic, wantMessages, messagesChan)

		hotelService := handlers.NewService(
			hotel_repo.NewRepo(dbrepo.NewPostgresRepo(db.DB)),
			producer.NewKafkaSender(tProducer.Producer, topic),
			parser_request.NewRequestParser(),
		)

		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		//act
		roomID, err := hotelService.CreateRoom(http.MethodPost, createBody, sync)
		require.Nil(t, err)

		err = hotelService.DeleteRoomWithAllReservations(http.MethodDelete, roomID, sync)
		require.Nil(t, err)

		for i := 0; i < wantMessages; i++ {
			<-messagesChan
		}
		w.Close()
		out, _ := io.ReadAll(r)
		os.Stdout = rescueStdout

		//assert
		require.Contains(t, string(out), `"method":"DELETE"`)
		require.Contains(t, string(out), createOutBody)
		require.Contains(t, string(out), deleteOutBody)
	})
}

func Test_CreateReservation(t *testing.T) {
	t.Parallel()
	var (
		createRoomBody   = []byte("{\"name\":\"Lux\", \"cost\":1000.0}")
		createResStrBody = "{\"start_date\":\"2023-09-15\", \"end_date\":\"2023-11-15\", \"room_id\":%v}"
		topic            = viper.GetString("TOPIC")
		outRoomBody      = "\"{\\\"name\\\":\\\"Lux\\\", \\\"cost\\\":1000.0}\""
		outResBody       = "\"{\\\"start_date\\\":\\\"2023-09-15\\\", \\\"end_date\\\":\\\"2023-11-15\\\", \\\"room_id\\\":%v}\""
		sync             = true
		wantMessages     = 2
	)
	t.Run("success", func(t *testing.T) {
		db.SetUp(t)
		defer db.TearDown()

		tProducer.SetUp(t)
		defer tProducer.TearDown()

		tConsumer.SetUp(t)
		defer tConsumer.TearDown()

		messagesChan := make(chan bool, wantMessages)
		tConsumer.Consumer.StartConsume(topic, wantMessages, messagesChan)

		hotelService := handlers.NewService(
			hotel_repo.NewRepo(dbrepo.NewPostgresRepo(db.DB)),
			producer.NewKafkaSender(tProducer.Producer, topic),
			parser_request.NewRequestParser(),
		)

		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		//act
		roomID, err := hotelService.CreateRoom(http.MethodPost, createRoomBody, sync)
		require.Nil(t, err)
		_, err = hotelService.CreateReservation(http.MethodPost, []byte(fmt.Sprintf(createResStrBody, roomID)), sync)
		require.Nil(t, err)

		for i := 0; i < wantMessages; i++ {
			<-messagesChan
		}
		w.Close()
		out, _ := io.ReadAll(r)
		os.Stdout = rescueStdout

		//assert
		require.Contains(t, string(out), outRoomBody)
		require.Contains(t, string(out), fmt.Sprintf(outResBody, roomID))
	})
}

func Test_UpdateReservationKafka(t *testing.T) {
	t.Parallel()
	var (
		createRoomBody      = []byte("{\"name\":\"Lux\", \"cost\":1000.0}")
		createResStrBody    = "{\"start_date\":\"2023-09-15\", \"end_date\":\"2023-11-15\", \"room_id\":%v}"
		updateResStrBody    = "{\"id\":%v, \"start_date\":\"2023-09-15\", \"end_date\":\"2023-11-15\", \"room_id\":%v}"
		topic               = viper.GetString("TOPIC")
		outRoomBody         = "\"{\\\"name\\\":\\\"Lux\\\", \\\"cost\\\":1000.0}\""
		outResBody          = "\"{\\\"start_date\\\":\\\"2023-09-15\\\", \\\"end_date\\\":\\\"2023-11-15\\\", \\\"room_id\\\":%v}\""
		outUpdateResStrBody = "{\\\"id\\\":%v, \\\"start_date\\\":\\\"2023-09-15\\\", \\\"end_date\\\":\\\"2023-11-15\\\", \\\"room_id\\\":%v}"
		sync                = true
		wantMessages        = 3
	)
	t.Run("success", func(t *testing.T) {
		db.SetUp(t)
		defer db.TearDown()

		tProducer.SetUp(t)
		defer tProducer.TearDown()

		tConsumer.SetUp(t)
		defer tConsumer.TearDown()

		messagesChan := make(chan bool, wantMessages)
		tConsumer.Consumer.StartConsume(topic, wantMessages, messagesChan)

		hotelService := handlers.NewService(
			hotel_repo.NewRepo(dbrepo.NewPostgresRepo(db.DB)),
			producer.NewKafkaSender(tProducer.Producer, topic),
			parser_request.NewRequestParser(),
		)

		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		//act
		roomID, err := hotelService.CreateRoom(http.MethodPost, createRoomBody, sync)
		require.Nil(t, err)
		resID, err := hotelService.CreateReservation(http.MethodPost, []byte(fmt.Sprintf(createResStrBody, roomID)), sync)
		require.Nil(t, err)
		err = hotelService.UpdateReservation(http.MethodPut, []byte(fmt.Sprintf(updateResStrBody, resID, roomID)), sync)
		require.Nil(t, err)

		for i := 0; i < wantMessages; i++ {
			<-messagesChan
		}
		w.Close()
		out, _ := io.ReadAll(r)
		os.Stdout = rescueStdout

		//assert
		require.Contains(t, string(out), `"method":"PUT"`)
		require.Contains(t, string(out), outRoomBody)
		require.Contains(t, string(out), fmt.Sprintf(outResBody, roomID))
		require.Contains(t, string(out), fmt.Sprintf(outUpdateResStrBody, resID, roomID))
	})
}

func Test_GetReservation(t *testing.T) {
	t.Parallel()
	var (
		createRoomBody   = []byte("{\"name\":\"Lux\", \"cost\":1000.0}")
		createResStrBody = "{\"start_date\":\"2023-09-15\", \"end_date\":\"2023-11-15\", \"room_id\":%v}"
		getOutBody       = "\"body\":\"\""
		topic            = viper.GetString("TOPIC")
		outRoomBody      = "\"{\\\"name\\\":\\\"Lux\\\", \\\"cost\\\":1000.0}\""
		outResBody       = "\"{\\\"start_date\\\":\\\"2023-09-15\\\", \\\"end_date\\\":\\\"2023-11-15\\\", \\\"room_id\\\":%v}\""
		sync             = true
		wantMessages     = 3
	)
	t.Run("success", func(t *testing.T) {
		db.SetUp(t)
		defer db.TearDown()

		tProducer.SetUp(t)
		defer tProducer.TearDown()

		tConsumer.SetUp(t)
		defer tConsumer.TearDown()

		messagesChan := make(chan bool, wantMessages)
		tConsumer.Consumer.StartConsume(topic, wantMessages, messagesChan)

		hotelService := handlers.NewService(
			hotel_repo.NewRepo(dbrepo.NewPostgresRepo(db.DB)),
			producer.NewKafkaSender(tProducer.Producer, topic),
			parser_request.NewRequestParser(),
		)

		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		//act
		roomID, err := hotelService.CreateRoom(http.MethodPost, createRoomBody, sync)
		require.Nil(t, err)
		resID, err := hotelService.CreateReservation(http.MethodPost, []byte(fmt.Sprintf(createResStrBody, roomID)), sync)
		require.Nil(t, err)
		_, err = hotelService.GetReservation(http.MethodGet, resID, sync)
		require.Nil(t, err)

		for i := 0; i < wantMessages; i++ {
			<-messagesChan
		}
		w.Close()
		out, _ := io.ReadAll(r)
		os.Stdout = rescueStdout

		//assert
		require.Contains(t, string(out), `"method":"GET"`)
		require.Contains(t, string(out), outRoomBody)
		require.Contains(t, string(out), fmt.Sprintf(outResBody, roomID))
		require.Contains(t, string(out), getOutBody)
	})
}

func Test_DeleteReservation(t *testing.T) {
	t.Parallel()
	var (
		createRoomBody   = []byte("{\"name\":\"Lux\", \"cost\":1000.0}")
		createResStrBody = "{\"start_date\":\"2023-09-15\", \"end_date\":\"2023-11-15\", \"room_id\":%v}"
		deleteOutBody    = "\"body\":\"\""
		topic            = viper.GetString("TOPIC")
		outRoomBody      = "\"{\\\"name\\\":\\\"Lux\\\", \\\"cost\\\":1000.0}\""
		outResBody       = "\"{\\\"start_date\\\":\\\"2023-09-15\\\", \\\"end_date\\\":\\\"2023-11-15\\\", \\\"room_id\\\":%v}\""
		sync             = true
		wantMessages     = 3
	)
	t.Run("success", func(t *testing.T) {
		db.SetUp(t)
		defer db.TearDown()

		tProducer.SetUp(t)
		defer tProducer.TearDown()

		tConsumer.SetUp(t)
		defer tConsumer.TearDown()

		messagesChan := make(chan bool, wantMessages)
		tConsumer.Consumer.StartConsume(topic, wantMessages, messagesChan)

		hotelService := handlers.NewService(
			hotel_repo.NewRepo(dbrepo.NewPostgresRepo(db.DB)),
			producer.NewKafkaSender(tProducer.Producer, topic),
			parser_request.NewRequestParser(),
		)

		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		//act
		roomID, err := hotelService.CreateRoom(http.MethodPost, createRoomBody, sync)
		require.Nil(t, err)
		resID, err := hotelService.CreateReservation(http.MethodPost, []byte(fmt.Sprintf(createResStrBody, roomID)), sync)
		require.Nil(t, err)
		err = hotelService.DeleteReservation(http.MethodDelete, resID, sync)
		require.Nil(t, err)

		for i := 0; i < wantMessages; i++ {
			<-messagesChan
		}
		w.Close()
		out, _ := io.ReadAll(r)
		os.Stdout = rescueStdout

		//assert
		require.Contains(t, string(out), `"method":"DELETE"`)
		require.Contains(t, string(out), outRoomBody)
		require.Contains(t, string(out), fmt.Sprintf(outResBody, roomID))
		require.Contains(t, string(out), deleteOutBody)
	})
}
