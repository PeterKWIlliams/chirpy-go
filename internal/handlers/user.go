package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/PeterKWIlliams/chirpy-go/internal/auth"
)

func (cfg *ApiCfg) HandlerPostUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid json")
		return
	}
	passwordHash, err := auth.HashPassword(params.Password)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "There was an error creating the user")
		return
	}
	user, error := cfg.Database.CreateUser(params.Email, passwordHash)
	if error != nil {
		RespondWithError(w, http.StatusBadRequest, "There was an error creating the chirp")
		return
	}
	strippedUser := struct {
		ID          int    `json:"id"`
		Email       string `json:"email"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}{
		ID:          user.ID,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	RespondWithJSON(w, http.StatusCreated, strippedUser)
}

func (cfg *ApiCfg) HandlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	tokenString, err := auth.ExtractBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	id, err := auth.VerifyJWT(tokenString, cfg.JWTSecret)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid json")
		return
	}

	passwordHash, err := auth.HashPassword(params.Password)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "There was an error creating the chirp")
		return
	}
	user, err := cfg.Database.UpdateUser(id, params.Email, passwordHash, tokenString)
	if err != nil {

		fmt.Println(err)
		RespondWithError(w, http.StatusUnauthorized, "Error updating user. Could not get user")
		return
	}

	strippedUser := struct {
		ID          int    `json:"id"`
		Email       string `json:"email"`
		IsChirpyRed bool   `json:"is_chirpy_red"`
	}{
		ID:          user.ID,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}
	RespondWithJSON(w, http.StatusOK, strippedUser)
}

func (cfg *ApiCfg) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}
	type response struct {
		ID           int    `json:"id"`
		Email        string `json:"email"`
		IsChirpyRed  bool   `json:"is_chirpy_red"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid json")
		return
	}
	user, err := cfg.Database.GetUserByEmail(params.Email)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "user not found")
		return
	}

	err = auth.VerifyPassword(user.Password, params.Password)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "invalid password")
		return
	}
	hourInSeconds := 60 * 60
	if params.ExpiresInSeconds == 0 || params.ExpiresInSeconds <= hourInSeconds {
		params.ExpiresInSeconds = hourInSeconds
	}
	token, err := auth.GenerateJWT(user.ID, params.ExpiresInSeconds, cfg.JWTSecret)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "error creating token")
		return
	}
	refreshToken, err := cfg.Database.CreateRefToken(user.ID)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "error creating token")
		return
	}
	RespondWithJSON(w, http.StatusOK, response{
		ID:           user.ID,
		Email:        user.Email,
		IsChirpyRed:  user.IsChirpyRed,
		Token:        token,
		RefreshToken: refreshToken,
	})
}
