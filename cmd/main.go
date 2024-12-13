package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/leonardonicola/golerplate/internal/config"
	"github.com/leonardonicola/golerplate/pkg/constants"
)

//	@title			Golerplate
//	@version		1.0
//	@description Boilerplate for Golang
//	@host			localhost:3000

// @contact.name	Autor
// @contact.url	https://github.com/leonardonicola
func main() {
	config.LoadEnv()

	dbConfig := config.NewDBConfig()
	db, err := config.NewDatabase(dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	router := config.NewRouter(db.Pool)

	srv := &http.Server{
		Addr:    constants.PORT,
		Handler: router.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Graceful shutdown...")

	// Give it 5 seconds to shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	select {
	case <-ctx.Done():
		log.Println("Timeout ended.")
	}
	log.Println("Server exiting...")
}
