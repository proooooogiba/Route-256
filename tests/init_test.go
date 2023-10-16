//go:build integration
// +build integration

package tests

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"homework-3/tests/postgres"
)

var (
	db *postgres.TDB
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}
	viper.AutomaticEnv()

	db = postgres.NewFromEnv()
}
