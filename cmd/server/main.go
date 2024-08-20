package main

import (
	"log"

	"github.com/PeterKWIlliams/chirpy-go/internal/server"
)

func main() {
	srv := server.NewServer()

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
