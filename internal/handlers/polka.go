package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/PeterKWIlliams/chirpy-go/internal/auth"
)

func (cfg *ApiCfg) HandlerUpgradeMembership(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.ExtractBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if apiKey != cfg.PolkaKey {
		RespondWithError(w, http.StatusUnauthorized, "Unauthorized")
	}
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		}
	}
	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if strings.Split(params.Event, ".")[1] != "upgraded" {
		RespondWithError(w, http.StatusNoContent, "invalid event")
		return
	}

	user, err := cfg.Database.GetUserById(params.Data.UserID)
	if err != nil {
		RespondWithError(w, http.StatusNotFound, "invalid user")
		return
	}
	if user.IsChirpyRed {
		RespondWithError(w, http.StatusNoContent, "user already has membership")
		return
	}
	user.IsChirpyRed = true

	err = cfg.Database.UpgradeUser(params.Data.UserID)
	if err != nil {
		RespondWithError(w, http.StatusNoContent, "could not upgrade user")
		return
	}
	RespondWithJSON(w, http.StatusNoContent, nil)
}
