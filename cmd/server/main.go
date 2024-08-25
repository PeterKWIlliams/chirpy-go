package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"

	"github.com/PeterKWIlliams/chirpy-go/internal/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading.env file")
	}
	projectRoot, _ := os.Getwd()
	path := "database.json"
	filePath := filepath.Join(projectRoot, path)
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	if *dbg {
		err := os.Remove(filePath)
		if err != nil {
			fmt.Println("Error removing database folder")
		}

	}
	srv := server.NewServer()

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
