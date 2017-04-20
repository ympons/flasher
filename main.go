package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ympons/flasher/db"
	"github.com/ympons/flasher/server"
)

func main() {
	port := getenv("PORT", "5900")
	basePath := getenv("BASE_PATH", "web")
	secretKey := getenv("SECRET_KEY_BASE", "8c18346e742aec88ddf68fc9f51e5e")
	dbSource := getenv("DATABASE_URL", "./db/cards-jwasham.db")

	dbInstance, err := db.Open("sqlite3", dbSource)
	if err != nil {
		log.Fatalf("Unable to connect to the DB: %v", err)
	}

	if err := dbInstance.InitDbSchemas(); err != nil {
		log.Fatalf("Initializing DB schemas failed: %v", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)

	srv := server.New(basePath, secretKey, dbInstance)

	// Block until a signal is received
	go func() {
		sig := <-c
		if err := srv.Close(); err != nil {
			log.Printf("Failed to close server: %v", err)
		}
		log.Printf("Exiting given signal: %v", sig)
		os.Exit(0)
	}()

	srv.Run(":" + port)
}

func getenv(k string, v string) string {
	if val := os.Getenv(k); val != "" {
		return val
	}
	os.Setenv(k, v)
	return v
}
