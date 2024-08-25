package database

import (
	"fmt"
	"log"
	"os"
)

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

func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	chirps := make([]Chirp, 0)

	data, err := db.loadDB()
	if err != nil {
		log.Println("Error loading db")
		return nil, err
	}
	for _, chirp := range data.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	data, err := db.loadDB()
	if err != nil {
		log.Println("Error loading db")
		return Chirp{}, err
	}
	chirp, ok := data.Chirps[id]
	if !ok {
		return Chirp{}, fmt.Errorf("Chirp with id %d does not exist", id)
	}
	return chirp, nil
}
