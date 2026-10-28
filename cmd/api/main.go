package main

import (
	"log"
	"net/http"
	"os"

	"libratrack/internal/app"
	"libratrack/internal/repository"
)

func main() {
	db, err := repository.NewDB()
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}

	if err := repository.RunMigrations(os.Getenv("DATABASE_URL"), "migrations"); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	addr := ":8080"
	port := os.Getenv("PORT")
	if port != "" {
		addr = ":" + port
	}

	log.Printf("listening on %s", addr)
	handler := app.NewRouterWithDB(db)
	err = http.ListenAndServe(addr, handler)
	if err != nil {
		log.Fatal(err)
	}
}
