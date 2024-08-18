package database

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

func NewDB(path string) (*DB, error) {
	projectRoot, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting project root path")
		return nil, err
	}

	filePath := filepath.Join(projectRoot, "example.txt")
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating file")
		return nil, err
	}

	defer file.Close()

	return &DB{mux: &sync.RWMutex{}, path: filePath}, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbData, err := db.loadDB()
	if err != nil {
		log.Println("Error loading db")
		return Chirp{}, err
	}

	_, err = os.Stat(db.path)
	if err != nil {
		log.Println("File does not exist")
		return Chirp{}, err
	}

	ID := len(dbData.Chirps) + 1
	chirp := Chirp{ID: ID, Body: body}
	dbData.Chirps[ID] = chirp

	err = db.writeDb(dbData)
	if err != nil {
		log.Println("Error writing to db")
		return Chirp{}, err
	}
	return chirp, nil
}

func (db *DB) ensureDB() {
	if _, err := os.Stat(db.path); os.IsNotExist(err) {
		dbStructure := DBStructure{Chirps: make(map[int]Chirp)}
		err = db.writeDb(dbStructure)
		if err != nil {
			log.Println("Error creating db")
		}
	}
}

func (db *DB) loadDB() (DBStructure, error) {
	data := DBStructure{}
	file, err := os.ReadFile(db.path)
	if err != nil {
		log.Println("Error reading file")
		return DBStructure{}, nil
	}
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Println("Error marshalling json")
		return DBStructure{}, nil
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
		log.Println("Error trying to write to db")
		return err
	}

	return nil
}
