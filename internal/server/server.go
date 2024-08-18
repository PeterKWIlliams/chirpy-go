package server

import (
	"net/http"
	"path/filepath"

	"github.com/PeterKWIlliams/chirpy-go/internal/handlers"
)

func NewServer() *http.Server {
	path, err := filepath.Abs("web/")
	apiCfg := &handlers.ApiCfg{
		FileserverHits: 0,
		PostChirpHits:  0,
	}
	if err != nil {
		panic(err)
	}
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(path)))
	mux := http.NewServeMux()
	mux.Handle("/app/*", apiCfg.MiddlewareMetricsInc(handler))

	mux.HandleFunc("POST /api/chirps", handlers.HandlerPost)
	mux.HandleFunc("GET /api/healthz", handlers.HandlerHealhtz)
	mux.HandleFunc("GET /admin/metrics/", apiCfg.HandlerMetrics)
	mux.HandleFunc("/api/reset", apiCfg.HandlerReset)

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}
