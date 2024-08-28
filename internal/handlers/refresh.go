package handlers

import (
	"net/http"

	"github.com/PeterKWIlliams/chirpy-go/internal/auth"
)

func (cfg *ApiCfg) HandlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.ExtractBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "issue while refreshing token")
		return
	}
	userId, err := cfg.Database.VerifyRefreshToken(token)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Could not verify token while refreshing")
		return
	}

	hourInSeconds := 60 * 60
	signedToken, err := auth.GenerateJWT(userId, hourInSeconds, cfg.JWTSecret)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Signed Token genration gone wrong")
		return
	}
	RespondWithJSON(w, http.StatusOK, map[string]string{"token": signedToken})
}

func (cfg *ApiCfg) HandlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.ExtractBearerToken(r.Header)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Token is invalid")
		return
	}
	err = cfg.Database.RevokeRefreshToken(token)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Error revoking token")
	}
	RespondWithJSON(w, http.StatusNoContent, nil)
}
