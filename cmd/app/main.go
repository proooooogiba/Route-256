package main

import (
	"context"
	"github.com/spf13/viper"
	"homework-3/internal/commands/invoker"
	"homework-3/internal/handlers"
	"homework-3/internal/pkg/db"
	"log"
	"net/http"
	"os"
)

func main() {
	commandInvoker, withCommand := invoker.NewCommandInvoker(os.Args)
	if withCommand {
		commandInvoker.ExecuteCommand()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	db, err := db.NewDB(ctx)
	if err != nil {
		log.Fatal("can't establish connection with database", err)
	}
	defer db.GetPool(ctx).Close()

	hotelRepo := handlers.NewRepo(db)

	srv := &http.Server{
		Addr:    viper.Get("PORT").(string),
		Handler: routes(hotelRepo),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
