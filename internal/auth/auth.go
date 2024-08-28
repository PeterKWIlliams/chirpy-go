package auth

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {
	data, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func ExtractBearerToken(header http.Header) (string, error) {
	tokenString := header.Get("Authorization")

	if tokenString == "" {
		return "", errors.New("")
	}
	strippedTokenString := tokenString[7:]

	return strippedTokenString, nil
}

func GenerateJWT(userId int, expiresIn int, JWTSecret string) (string, error) {
	signKey := []byte(JWTSecret)

	issuedAt := jwt.NumericDate{Time: time.Now().UTC()}
	expiresAt := jwt.NumericDate{Time: time.Now().UTC().Add(time.Second * time.Duration(expiresIn))}
	subject := strconv.Itoa(userId)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  &issuedAt,
		ExpiresAt: &expiresAt,
		Subject:   subject,
	})
	signedToken, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

type TokenClaims struct {
	jwt.RegisteredClaims
}

func VerifyJWT(tokenString string, JWTSecret string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})
	if err != nil {
		return 0, err
	}
	id, err := strconv.Atoi(token.Claims.(*TokenClaims).Subject)
	if err != nil {
		return 0, err
	}
	return id, err
}

func VerifyPassword(password []byte, passwordHash string) error {
	return bcrypt.CompareHashAndPassword(password, []byte(passwordHash))
}
