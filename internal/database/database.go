package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Password    []byte `json:"password"`
	IsChirpyRed bool   `json:"isChirpyRed"`
}

type RefreshToken struct {
	UserId    int       `json:"userId"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type DBStructure struct {
	Chirps        map[int]Chirp           `json:"chirps"`
	Users         map[int]User            `json:"users"`
	UserEmail     map[string]int          `json:"userEmail"`
	RefreshTokens map[string]RefreshToken `json:"refreshTokens"`
}

func NewDB(path string) (*DB, error) {
	projectRoot, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting project root path")
		return nil, err
	}
	filePath := filepath.Join(projectRoot, path)

	db := &DB{mux: &sync.RWMutex{}, path: filePath}

	err = db.ensureDB()
	if err != nil {
		return &DB{}, nil
	}
	return db, nil
}

func (db *DB) ensureDB() error {
	if _, err := os.Stat(db.path); os.IsNotExist(err) {
		dbStructure := DBStructure{
			Chirps:        make(map[int]Chirp),
			Users:         make(map[int]User),
			UserEmail:     make(map[string]int),
			RefreshTokens: make(map[string]RefreshToken),
		}
		err = db.writeDb(dbStructure)
		if err != nil {
			log.Println("Error creating db")
			return err
		}
	}
	return nil
}

func (db *DB) loadDB() (DBStructure, error) {
	data := DBStructure{}
	file, err := os.ReadFile(db.path)
	if err != nil {
		log.Println("Error reading file while loading db")
		return DBStructure{}, err
	}
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Println("Error marshalling json")
		return DBStructure{}, err
	}
	return data, nil
}

func (db *DB) writeDb(dbStructure DBStructure) error {
	data, err := json.Marshal(dbStructure)
	if err != nil {
		log.Println("Error marshalling json while writing to db")
		return err
	}
	err = os.WriteFile(db.path, data, os.FileMode(0644))
	if err != nil {
		log.Println("Error trying to write to db while writing to db")
		return err
	}

	return nil
}
