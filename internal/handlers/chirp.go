package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/PeterKWIlliams/chirpy-go/internal/auth"
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

func (cfg *ApiCfg) HandlerPost(w http.ResponseWriter, r *http.Request) {
	token, err := auth.ExtractBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	id, err := auth.VerifyJWT(token, cfg.JWTSecret)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized could not verify token")
		return
	}
	type parameters struct {
		Body string `json:"body"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if len(params.Body) > 140 {
		RespondWithError(w, http.StatusBadRequest, "body is too long")
		return
	}

	chirp, error := cfg.Database.CreateChirp(removeProfanity(params.Body), id)
	if error != nil {
		RespondWithError(w, http.StatusBadRequest, "There was an error creating the chirp")
		return
	}
	RespondWithJSON(w, 201, chirp)
}

func (cfg *ApiCfg) HandlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	token, err := auth.ExtractBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	authorId, err := auth.VerifyJWT(token, cfg.JWTSecret)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized could not verify token")
		return
	}
	chirpIdString := r.PathValue("chirpId")
	id, err := strconv.Atoi(chirpIdString)
	fmt.Println("this is the id", chirpIdString)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid id")
		return
	}
	chirp, err := cfg.Database.GetChirp(id)
	if err != nil {
		RespondWithError(w, 404, "there was an error getting the chirp")
		return
	}
	if chirp.AuthorId != authorId {
		RespondWithError(w, http.StatusForbidden, "Unauthorized not allowed to delete chirp")
		return
	}
	RespondWithJSON(w, http.StatusNoContent, nil)
}

func (cfg *ApiCfg) HandlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpIdString := r.PathValue("chirpId")
	id, err := strconv.Atoi(chirpIdString)
	fmt.Println("this is the id", chirpIdString)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid id")
		return
	}
	chirp, err := cfg.Database.GetChirp(id)
	if err != nil {
		RespondWithError(w, 404, "there was an error getting the chirp")
		return
	}
	RespondWithJSON(w, http.StatusOK, chirp)
}

func (cfg *ApiCfg) HandlerGetChirps(w http.ResponseWriter, r *http.Request) {
	idString := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")

	chirps, err := cfg.Database.GetChirps()
	if err != nil {
		fmt.Println("Error getting chirps")
		RespondWithError(w, http.StatusInternalServerError, "There was an error getting the chirps")
		return
	}

	if idString != "" {
		id, err := strconv.Atoi(idString)
		if err != nil {
			RespondWithError(w, http.StatusBadRequest, "invalid id")
			return
		}
		chirps, err = cfg.Database.GetUserChirps(id)
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "There was an error getting the chirps")
			return
		}
	}

	fmt.Println(sortOrder, "this is the sort order")
	if sortOrder == "asc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID < chirps[j].ID
		})
	} else if sortOrder == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].ID > chirps[j].ID
		})
	}

	RespondWithJSON(w, http.StatusOK, chirps)
}
