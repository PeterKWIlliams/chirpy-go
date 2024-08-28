package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/PeterKWIlliams/chirpy-go/internal/database"
)

type ApiCfg struct {
	FileserverHits int
	Database       *database.DB
	JWTSecret      string
	PolkaKey       string
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	type errorResponse struct {
		Error string `json:"error"`
	}

	RespondWithJSON(w, code, errorResponse{Error: msg})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Add("content-type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}
