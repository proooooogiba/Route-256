package main

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"homework-3/internal/grpc_handlers"
	"homework-3/internal/pkg/db"
	"homework-3/internal/pkg/domain/hotel_repo"
	"homework-3/internal/pkg/grpc_server"
	"homework-3/internal/pkg/logger"
	"homework-3/internal/pkg/pb"
	"homework-3/internal/pkg/repository/dbrepo"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	logger.SetGlobal(
		zapLogger.With(zap.String("server", "gRPC")),
	)

	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 1 * time.Second,
		},
	}
	tracer, closer, err := cfg.New(
		"hotel-service",
	)
	if err != nil {
		logger.Errorf(ctx, "can't create traces", err)
	}
	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)

	db, err := db.NewDB(ctx)
	if err != nil {
		logger.Errorf(ctx, "can't establish connection with database", err)
		os.Exit(1)
	}
	defer db.GetPool(ctx).Close()
	logger.Infof(ctx, "connection with database is established")

	hotelService := grpc_handlers.NewService(
		hotel_repo.NewRepo(dbrepo.NewPostgresRepo(db)),
	)

	server := grpc.NewServer()
	reflection.Register(server)

	pb.RegisterHotelServiceServer(server, grpc_server.NewImplementation(hotelService))

	listener, err := net.Listen("tcp", viper.GetString("GRPC_PORT"))
	if err != nil {
		log.Fatal(err)
	}

	logger.Infof(ctx, "starting server on %s", listener.Addr().String())
	if err := server.Serve(listener); err != nil {
		logger.Errorf(ctx, "fail while serve gRPC server", err)
		os.Exit(1)
	}
}
