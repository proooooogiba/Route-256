package db

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewDB(ctx context.Context) (*Database, error) {
	pool, err := pgxpool.Connect(ctx, generateDsn())
	if err != nil {
		return nil, err
	}
	log.Println("Connected to database...")
	return newDatabase(pool), nil
}

func generateDsn() string {
	return fmt.Sprintf("host=%s port=%s user=test password=%s dbname=%s sslmode=disable",
		viper.Get("HOST"),
		viper.Get("POSTGRES_PORT"),
		//viper.Get("USER"),
		viper.Get("PASSWORD"),
		viper.Get("DBNAME"))
}
