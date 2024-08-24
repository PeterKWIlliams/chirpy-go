package handlers

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
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

	user, error := cfg.Database.CreateUser(params.Email, params.Password)
	if error != nil {
		RespondWithError(w, http.StatusBadRequest, "There was an error creating the chirp")
		return
	}
	strippedUser := struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}{ID: user.ID, Email: user.Email}

	RespondWithJSON(w, http.StatusCreated, strippedUser)
}

func (cfg *ApiCfg) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
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
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(params.Password)); err != nil {
		RespondWithError(w, http.StatusUnauthorized, "invalid password")
		return
	}
	RespondWithJSON(w, http.StatusOK, user)
}
