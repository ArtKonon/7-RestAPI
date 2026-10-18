package main

import (
	"log"
	"net/http"
	"os"

	"homework/internal/app"
)

func main() {
	addr := ":8080"
	port := os.Getenv("PORT")
	if port != "" {
		addr = ":" + port
	}

	log.Printf("listening on %s", addr)
	handler := app.NewRouter()
	err := http.ListenAndServe(addr, handler)
	if err != nil {
		log.Fatal(err)
	}
}
