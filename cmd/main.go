package main

import (
	"context"
	// "errors"
	"log"
	// "net/http"
	// "os"
	// "os/signal"
	// "syscall"
	"time"

	"github.com/ebriussenex/shrt/internal/db"
	"github.com/ebriussenex/shrt/internal/server"
	"github.com/ebriussenex/shrt/internal/shorten"
	shorten_store "github.com/ebriussenex/shrt/internal/storage"

	config "github.com/ebriussenex/shrt/internal/configs"
)

func main() {
	dbCtx, dbCancel := context.WithTimeout(context.Background(), 1 * time.Second)
	defer dbCancel()

	mongoClient, err := db.Connect(dbCtx, config.Get().DB.DSN)
	if err != nil {
		log.Fatal(err)
	}


	shortenStorage := shorten_store.NewMongoDB(mongoClient.Client())
	shortenService := shorten.NewService(shortenStorage)
	srv := server.New(shortenService)
	
	server.StartServerWithGracefulShutdown(srv)

	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	// go func() {
	// 	if err := http.ListenAndServe(config.Get().ListenAddr(), srv); !errors.Is(err, http.ErrServerClosed) {
	// 		log.Fatalf("error running server: %v", err)
	// 	}
	// }()

	// log.Printf("started server\n")
	// <-quit

	// shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer shutdownCancel()

	// if err := srv.Shutdown(shutdownCtx); err != nil {
	// 	log.Fatalf("error closing server: %v", err)
	// }

	// log.Println("server stopped")
}