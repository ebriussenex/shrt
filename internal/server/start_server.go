package server

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"

	config "github.com/ebriussenex/shrt/internal/configs"
)

func StartServerWithGracefulShutdown(a *Server) {
	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := a.f.Shutdown(); err != nil {
			log.Printf("Server is not shutting down! Reason: %v", err)
		}
		close(idleConnsClosed)
	}()
	cfg := config.Get()
	fiberConnURL := fmt.Sprintf("%s:%s", cfg.Host, strconv.Itoa(cfg.Port))

	if err := a.f.Listen(fiberConnURL); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}

	<-idleConnsClosed
}
