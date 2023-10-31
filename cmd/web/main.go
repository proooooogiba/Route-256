package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"homework-3/infrastructure/kafka"
	"homework-3/internal/consumer"
	"homework-3/internal/handlers"
	"homework-3/internal/pkg/bussiness_logic/hotel_repo"
	"homework-3/internal/pkg/db"
	"homework-3/internal/pkg/parser/parser_request"
	"homework-3/internal/pkg/repository/dbrepo"
	"homework-3/internal/producer"
	"log"
	"net/http"
	"strings"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	db, err := db.NewDB(ctx)
	if err != nil {
		log.Fatal("can't establish connection with database", err)
	}
	defer db.GetPool(ctx).Close()

	brokersValues := viper.GetString("BROKERS")
	brokers := strings.Split(brokersValues, ",")

	kafkaProducer, err := kafka.NewProducer(brokers)
	if err != nil {
		log.Fatal(err)
	}
	defer kafkaProducer.Close()

	kafkaConsumer, err := consumer.NewConsumerService(brokers)
	if err != nil {
		log.Fatal(err)
	}
	defer kafkaConsumer.Close()

	topic := viper.Get("TOPIC").(string)
	kafkaConsumer.StartConsume(topic, -1, make(chan bool))

	hotelService := handlers.NewService(
		hotel_repo.NewRepo(dbrepo.NewPostgresRepo(db)),
		producer.NewKafkaSender(kafkaProducer, topic),
		parser_request.NewRequestParser(),
	)

	srv := &http.Server{
		Addr:    viper.Get("PORT").(string),
		Handler: routes(hotelService),
	}

	err = srv.ListenAndServe()
	fmt.Println(err)
}
