package database

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

func (db *DB) VerifyRefreshToken(token string) (int, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	data, err := db.loadDB()
	if err != nil {
		return 0, err
	}

	if _, ok := data.RefreshTokens[token]; !ok {
		return 0, errors.New("invalid token")
	}
	if time.Now().After(data.RefreshTokens[token].ExpiresAt) {
		return 0, errors.New("token expired")
	}

	userId := data.RefreshTokens[token].UserId
	return userId, nil
}

func (db *DB) CreateRefToken(userId int) (string, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := db.loadDB()
	if err != nil {
		return "", err
	}

	refreshTokenString, err := GenerateRefToken()
	if err != nil {
		return "", err
	}
	refreshToken := RefreshToken{
		UserId:    userId,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
	}

	data.RefreshTokens[refreshTokenString] = refreshToken
	err = db.writeDb(data)
	if err != nil {
		return "", err
	}
	return refreshTokenString, nil
}

func GenerateRefToken() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	refToken := hex.EncodeToString(randomBytes)
	return refToken, nil
}

func (db *DB) RevokeRefreshToken(token string) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := db.loadDB()
	if err != nil {
		return err
	}

	if _, ok := data.RefreshTokens[token]; !ok {
		return errors.New("Could not find token to revoke")
	}
	delete(data.RefreshTokens, token)
	err = db.writeDb(data)
	if err != nil {
		return err
	}
	return nil
}
