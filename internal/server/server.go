package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/PeterKWIlliams/chirpy-go/internal/database"
	"github.com/PeterKWIlliams/chirpy-go/internal/handlers"
)

func NewServer() *http.Server {
	path, _ := filepath.Abs("web/")

	databases, err := database.NewDB("database.json")
	if err != nil {
		fmt.Println("there was an eror creating the database")
	}
	apiCfg := &handlers.ApiCfg{
		FileserverHits: 0,
		Database:       databases,
		JWTSecret:      os.Getenv("JWT_TOKEN"),
	}
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(path)))
	mux := http.NewServeMux()
	mux.Handle("/app/*", apiCfg.MiddlewareMetricsInc(handler))

	mux.HandleFunc("POST /api/chirps", apiCfg.HandlerPost)
	mux.HandleFunc("GET /api/chirps", apiCfg.HandlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpId}", apiCfg.GetChirp)

	mux.HandleFunc("POST /api/login", apiCfg.HandlerLogin)
	mux.HandleFunc("POST /api/users", apiCfg.HandlerPostUser)
	mux.HandleFunc("PUT /api/users", apiCfg.HandlerUpdateUser)

	mux.HandleFunc("/api/reset", apiCfg.HandlerReset)
	mux.HandleFunc("GET /api/healthz", handlers.HandlerHealhtz)
	mux.HandleFunc("GET /admin/metrics/", apiCfg.HandlerMetrics)
	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}
