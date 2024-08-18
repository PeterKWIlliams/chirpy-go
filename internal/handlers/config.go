package handlers

import (
	"fmt"
	"net/http"
)

type ApiCfg struct {
	FileserverHits int
	PostChirpHits  int
}

func (cfg *ApiCfg) HandlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.FileserverHits = 0
	fmt.Println("server hit:", cfg.FileserverHits)
}

func (cfg *ApiCfg) HandlerPostChirp(w http.ResponseWriter, r *http.Request) {
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
