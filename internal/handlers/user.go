package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "There was an error creating the chirp")
		return
	}
	user, error := cfg.Database.CreateUser(params.Email, passwordHash)
	if error != nil {
		RespondWithError(w, http.StatusBadRequest, "There was an error creating the chirp")
		return
	}
	strippedUser := struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}{
		ID:    user.ID,
		Email: user.Email,
	}

	RespondWithJSON(w, http.StatusCreated, strippedUser)
}

type TokenClaims struct {
	jwt.RegisteredClaims
}

func (cfg *ApiCfg) HandlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		RespondWithError(w, http.StatusUnauthorized, "no token present:invalid token")
		return
	}
	tokenString = tokenString[7:]
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "invalid token heyayyayya")
		return
	}
	id, err := strconv.Atoi(token.Claims.(*TokenClaims).Subject)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "could not get id from subject")
		return
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid json")
		return
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, "There was an error creating the chirp")
		return
	}
	user, err := cfg.Database.UpdateUser(id, params.Email, passwordHash)
	fmt.Println("this is the password", params.Password)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Error updating user. Could not get user")
		return
	}

	strippedUser := struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}{
		ID:    user.ID,
		Email: user.Email,
	}
	RespondWithJSON(w, http.StatusOK, strippedUser)
}

func (cfg *ApiCfg) HandlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
	}
	type response struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		Token string `json:"token"`
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
	if params.ExpiresInSeconds == nil {
		dayInSeconds := 24 * 60 * 60
		params.ExpiresInSeconds = &dayInSeconds

	}
	issuedAt := jwt.NumericDate{Time: time.Now().UTC()}
	expiresAt := jwt.NumericDate{Time: time.Now().UTC().Add(time.Second * time.Duration(*params.ExpiresInSeconds))}
	subject := strconv.Itoa(user.ID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  &issuedAt,
		ExpiresAt: &expiresAt,
		Subject:   subject,
	})
	secret := []byte(cfg.JWTSecret)
	signedToken, err := token.SignedString(secret)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "error creating token")
		return
	}
	RespondWithJSON(w, http.StatusOK, response{
		ID:    user.ID,
		Email: user.Email,
		Token: signedToken,
	})
}
