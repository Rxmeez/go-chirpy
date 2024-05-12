package database

import (
	"errors"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var ErrorUserNotFound = errors.New("User Not Found")
var ErrorDuplicatedUser = errors.New("User has already been created")

type User struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
	RefreshToken
}

type RefreshToken struct {
	Token     string `json:"refresh_token"`
	expiresIn time.Duration
}

func (db *DB) CreateUser(email string, password string) (User, error) {

	dbStructure, err := db.loadDB()
	if err != nil && !errors.Is(err, ErrorEmptyFile) {
		log.Fatal(err)
		return User{}, err
	}

	if errors.Is(err, ErrorEmptyFile) {
		dbStructure = DBStructure{
			Chirps: make(map[int]Chirp),
			Users:  make(map[int]User),
		}
	}

	if _, err := dbStructure.findUserByEmail(email); err == nil {
		return User{}, errors.New("User has already been created")
	}

	userId := len(dbStructure.Users) + 1
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	user := User{Id: userId, Email: email, Password: string(hashPassword), IsChirpyRed: false}

	dbStructure.Users[userId] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		log.Fatal(err)
		return User{}, err
	}

	return user, nil
}

func (db *DB) UpdateUser(userId int, newEmail string, newPassword string) (User, error) {

	dbStructure, err := db.loadDB()
	if err != nil && !errors.Is(err, ErrorEmptyFile) {
		log.Fatal(err)
		return User{}, err
	}

	if errors.Is(err, ErrorEmptyFile) {
		dbStructure = DBStructure{
			Chirps: make(map[int]Chirp),
			Users:  make(map[int]User),
		}
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	user, ok := dbStructure.Users[userId]
	if !ok {
		return User{}, errors.New("Unable to find user")
	}

	user.Email = newEmail
	user.Password = string(hashPassword)

	dbStructure.Users[userId] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		log.Fatal(err)
		return User{}, err
	}

	return user, nil
}

func (db *DB) UpgradeUser(userId int) (User, error) {

	dbStructure, err := db.loadDB()
	if err != nil && !errors.Is(err, ErrorEmptyFile) {
		log.Fatal(err)
		return User{}, err
	}

	if errors.Is(err, ErrorEmptyFile) {
		dbStructure = DBStructure{
			Chirps: make(map[int]Chirp),
			Users:  make(map[int]User),
		}
	}

	user, ok := dbStructure.Users[userId]
	if !ok {
		return User{}, errors.New("Unable to find user")
	}

	user.IsChirpyRed = true

	dbStructure.Users[userId] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		log.Fatal(err)
		return User{}, err
	}

	return user, nil
}

func (d *DBStructure) findUserByEmail(email string) (User, error) {
	for _, user := range d.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return User{}, ErrorUserNotFound
}

func (db *DB) Login(email string, password string) (User, error) {

	dbStructure, err := db.loadDB()
	if err != nil && !errors.Is(err, ErrorEmptyFile) {
		log.Fatal(err)
		return User{}, err
	}

	if errors.Is(err, ErrorEmptyFile) {
		dbStructure = DBStructure{
			Chirps: make(map[int]Chirp),
			Users:  make(map[int]User),
		}
	}

	user, err := dbStructure.findUserByEmail(email)
	if err != nil {
		return User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (db *DB) StoreRefreshToken(userId int, refreshTokenString string, expireIn time.Duration) error {

	dbStructure, err := db.loadDB()
	if err != nil && !errors.Is(err, ErrorEmptyFile) {
		log.Fatal(err)
		return err
	}

	if errors.Is(err, ErrorEmptyFile) {
		dbStructure = DBStructure{
			Chirps: make(map[int]Chirp),
			Users:  make(map[int]User),
		}
	}

	user, ok := dbStructure.Users[userId]
	if !ok {
		return errors.New("Couldn't find the user")
	}

	user.RefreshToken.Token = refreshTokenString
	user.RefreshToken.expiresIn = expireIn

	dbStructure.Users[userId] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (db *DB) ValidateRefreshToken(refreshToken string) (int, error) {
	dbStructure, err := db.loadDB()
	if err != nil && !errors.Is(err, ErrorEmptyFile) {
		log.Fatal(err)
		return 0, err
	}

	if errors.Is(err, ErrorEmptyFile) {
		dbStructure = DBStructure{
			Chirps: make(map[int]Chirp),
			Users:  make(map[int]User),
		}
	}

	for _, user := range dbStructure.Users {
		if user.RefreshToken.Token == refreshToken {
			return user.Id, nil
		}
	}

	return 0, errors.New("Unable to Validate Refresh Token")

}

func (db *DB) RevokeRefreshToken(refreshToken string) error {
	dbStructure, err := db.loadDB()
	if err != nil && !errors.Is(err, ErrorEmptyFile) {
		log.Fatal(err)
		return err
	}

	if errors.Is(err, ErrorEmptyFile) {
		dbStructure = DBStructure{
			Chirps: make(map[int]Chirp),
			Users:  make(map[int]User),
		}
	}

	for _, user := range dbStructure.Users {
		if user.RefreshToken.Token == refreshToken {
			user.RefreshToken.Token = ""

			dbStructure.Users[user.Id] = user
			err = db.writeDB(dbStructure)
			if err != nil {
				log.Fatal(err)
				return err
			}
			return nil
		}
	}

	return errors.New("Unable to find token to revoke")

}
