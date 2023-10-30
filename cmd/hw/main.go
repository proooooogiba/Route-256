package main

import (
	"context"
	"fmt"
	"gitlab.ozon.dev/go/classroom-9/students/homework-7/internal/controller"
	"gitlab.ozon.dev/go/classroom-9/students/homework-7/internal/datasource/cache"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	cacheDB, err := cache.NewDatabaseWithCacheClient(
		fmt.Sprintf("./db/%s.json", time.Now().Format("2006-01-02 15:04:05")),
	)
	if err != nil {
		log.Fatal(err)
	}

	client := controller.NewClient(cacheDB) // TODO

	var bestUser = "best user, expired after 30 seconds"

	// Создаём запись
	err = client.Set(ctx, "user:12345:profile", bestUser, 1)
	if err != nil {
		panic(err)
	}

	// Получаем запись из кэша
	got, err := client.Get(ctx, "user:12345:profile")
	if err != nil {
		panic(err)
	}

	fmt.Println(got)

	if got != bestUser {
		panic("invalid value")
	}

	select {
	case <-time.After(2 * time.Second):
	case <-ctx.Done():
	}

	// Получаем запись из базы данных и обновляем кэш
	gotAgain, err := client.Get(ctx, "user:12345:profile")
	if err != nil {
		panic(err)
	}

	fmt.Println(gotAgain)

	if gotAgain != bestUser {
		panic("invalid value")
	}
}
