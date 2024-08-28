package database

import (
	"errors"
	"log"
)

func (db *DB) CreateUser(email string, password []byte) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()

	dbData, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	if _, ok := dbData.UserEmail[email]; ok {
		return User{}, errors.New("user already exists")
	}
	ID := len(dbData.Users) + 1
	user := User{ID: ID, Email: email, Password: password, IsChirpyRed: false}
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

func (db *DB) UpgradeUser(id int) error {
	db.mux.Lock()
	defer db.mux.Unlock()
	dbData, err := db.loadDB()
	if err != nil {
		return errors.New("error loading db")
	}
	user, ok := dbData.Users[id]
	if !ok {
		return errors.New("no user with that id")
	}
	user.IsChirpyRed = true
	dbData.Users[id] = user
	err = db.writeDb(dbData)
	if err != nil {
		log.Println("Error writing to db when upgrading user")
		return errors.New("error writing to db")
	}
	return nil
}

func (db *DB) UpdateUser(id int, email string, password []byte, tokenString string) (User, error) {
	db.mux.Lock()
	defer db.mux.Unlock()
	dbData, err := db.loadDB()
	if err != nil {
		log.Println("Error when updating user")
		return User{}, err
	}
	user, ok := dbData.Users[id]
	if !ok {
		return User{}, errors.New("no user with that id")
	}
	delete(dbData.UserEmail, user.Email)
	user.Email = email
	user.Password = password

	dbData.Users[id] = user

	dbData.UserEmail[email] = id

	err = db.writeDb(dbData)
	if err != nil {
		log.Println("Error writing to db when updating user")
		return User{}, err
	}
	return user, nil
}

func (db *DB) GetUserById(id int) (User, error) {
	dbData, err := db.loadDB()
	if err != nil {
		log.Println("Error when getting user by id")
		return User{}, err
	}
	return dbData.Users[id], nil
}
