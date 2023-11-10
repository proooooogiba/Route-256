package main

import (
	"context"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"homework-3/internal/grpc_handlers"
	"homework-3/internal/pkg/db"
	"homework-3/internal/pkg/domain/hotel_repo"
	"homework-3/internal/pkg/grpc_server"
	"homework-3/internal/pkg/pb"
	"homework-3/internal/pkg/repository/dbrepo"
	"log"
	"net"
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

	//brokersValues := viper.GetString("BROKERS")
	//brokers := strings.Split(brokersValues, ",")
	//
	//kafkaProducer, err := kafka.NewProducer(brokers)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer kafkaProducer.Close()
	//
	//kafkaConsumer, err := consumer.NewConsumerService(brokers)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer kafkaConsumer.Close()
	//
	//topic := viper.Get("TOPIC").(string)
	//kafkaConsumer.StartConsume(topic, -1, make(chan bool))

	hotelService := grpc_handlers.NewService(
		hotel_repo.NewRepo(dbrepo.NewPostgresRepo(db)),
	)

	//hotelService := handlers.NewService(
	//	hotel_repo.NewRepo(dbrepo.NewPostgresRepo(db)),
	//	producer.NewKafkaSender(kafkaProducer, topic),
	//	parser_request.NewRequestParser(),
	//)

	//srv := &http.Server{
	//	Addr:    viper.Get("PORT").(string),
	//	Handler: http_server.Routes(hotelService),
	//}
	//
	//err = srv.ListenAndServe()
	//fmt.Println(err)

	// транспорт для grpc запросов
	server := grpc.NewServer()
	reflection.Register(server)

	pb.RegisterHotelServiceServer(server, grpc_server.NewImplementation(hotelService))
	//pb.RegisterStudentServiceServer(server, grpc_server.NewImplementation(hotelService))

	listener, err := net.Listen("tcp", viper.GetString("GRPC_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("starting server on %s", listener.Addr().String())
	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
