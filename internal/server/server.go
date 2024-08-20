package server

import (
	"fmt"
	"net/http"
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
	}
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(path)))
	mux := http.NewServeMux()
	mux.Handle("/app/*", apiCfg.MiddlewareMetricsInc(handler))
	mux.HandleFunc("POST /api/chirps", apiCfg.HandlerPost)
	mux.HandleFunc("GET /api/chirps", apiCfg.HandlerGetChirps)
	mux.HandleFunc("GET /api/healthz", handlers.HandlerHealhtz)
	mux.HandleFunc("GET /admin/metrics/", apiCfg.HandlerMetrics)
	mux.HandleFunc("/api/reset", apiCfg.HandlerReset)

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}
