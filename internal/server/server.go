package server

import (
	"fmt"
	"net/http"
	"path/filepath"
)

type apiCfg struct{ fileserverHits int }

func (cfg *apiCfg) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits++
		fmt.Printf("server hit: %d\n", cfg.fileserverHits)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiCfg) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	fmt.Println("server hit:", cfg.fileserverHits)
}

func (cfg *apiCfg) handlerMetrics() func(w http.ResponseWriter, r *http.Request) {
	html := `
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
`
	return func(w http.ResponseWriter, r *http.Request) {
		html = fmt.Sprintf(html, cfg.fileserverHits)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
	}
}

func handlerHealhtz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err := w.Write([]byte("OK"))
	if err != nil {
		panic(err)
	}
}

func NewServer() *http.Server {
	path, err := filepath.Abs("web/")

	cfg := &apiCfg{fileserverHits: 0}
	if err != nil {
		panic(err)
	}
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(path)))
	mux := http.NewServeMux()
	mux.Handle("/app/*", cfg.middlewareMetricsInc(handler))

	mux.HandleFunc("GET /api/healthz", handlerHealhtz)
	mux.HandleFunc("GET /admin/metrics/", cfg.handlerMetrics())
	mux.HandleFunc("/api/reset", cfg.handlerReset)

	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}
