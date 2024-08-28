package handlers

import (
	"fmt"
	"net/http"
)

func (cfg *ApiCfg) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits++
		fmt.Printf("server hit: %d\n", cfg.FileserverHits)
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiCfg) HandlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.FileserverHits = 0
	fmt.Println("server hit:", cfg.FileserverHits)
}

func (cfg *ApiCfg) HandlerMetrics(w http.ResponseWriter, r *http.Request) {
	html := `
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>
`
	html = fmt.Sprintf(html, cfg.FileserverHits)
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
