package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func removeProfanity(text string) string {
	words := strings.Fields(text)
	for i, word := range words {
		word = strings.ToLower(word)
		if word == "kerfuffle" || word == "sharbert" || word == "fornax" {
			words[i] = "****"
		}
	}
	text = strings.Join(words, " ")
	return text
}

func HandlerPost(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	length := len(params.Body)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if length > 140 {
		RespondWithError(w, http.StatusBadRequest, "body is too long")
		return
	}

	payload := struct {
		CleanedBody string `json:"cleaned_body"`
	}{
		CleanedBody: removeProfanity(params.Body),
	}
	RespondWithJSON(w, http.StatusOK, payload)
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
