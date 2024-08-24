package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

type DBStructure struct {
	Chirps    map[int]Chirp  `json:"chirps"`
	Users     map[int]User   `json:"users"`
	UserEmail map[string]int `json:"userEmail"`
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

func (db *DB) CreateUser(email string, password string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbData, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	if _, ok := dbData.UserEmail[email]; ok {
		return User{}, errors.New("user already exists")
	}

	ID := len(dbData.Users) + 1
	user := User{ID: ID, Email: email, Password: passwordHash}
	dbData.Users[ID] = user
	dbData.UserEmail[email] = ID

	err = db.writeDb(dbData)
	if err != nil {
		log.Println("Error writing to db when creating user")
		return User{}, err
	}
	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbData, err := db.loadDB()
	if err != nil {
		log.Println("Error when getting user by email")
		return User{}, err
	}
	id, ok := dbData.UserEmail[email]
	if !ok {
		return User{}, errors.New("no user with that email")
	}
	user, ok := dbData.Users[id]
	if !ok {
		return User{}, errors.New("no user with that ")
	}

	return user, nil
}

func (db *DB) getUserById(id int) (User, error) {
	dbData, err := db.loadDB()
	if err != nil {
		log.Println("Error when getting user by id")
		return User{}, err
	}
	return dbData.Users[id], nil
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

func (db *DB) ensureDB() error {
	if _, err := os.Stat(db.path); os.IsNotExist(err) {
		dbStructure := DBStructure{Chirps: make(map[int]Chirp), Users: make(map[int]User), UserEmail: make(map[string]int)}
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
