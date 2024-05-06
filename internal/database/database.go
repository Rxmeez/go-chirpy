package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"sync"
)

var ErrorEmptyFile = errors.New("EmptyFile")

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DB struct {
	path string
	mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

func NewDB(path string) (*DB, error) {
	if _, err := os.Stat(path); err == nil {
		err := os.Remove(path)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	}
	file, err := os.Create(path)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &DB{
		path: path,
		mux:  &sync.RWMutex{},
	}, nil

}

func (db *DB) loadDB() (DBStructure, error) {
	content, err := os.ReadFile(db.path)
	if len(content) == 0 {
		return DBStructure{}, ErrorEmptyFile
	}
	if err != nil {
		log.Fatal(err)
		return DBStructure{}, err
	}

	chirps := DBStructure{}

	err = json.Unmarshal(content, &chirps)
	if err != nil {
		log.Fatal(err)
		return DBStructure{}, err
	}

	return chirps, nil

}

func (db *DB) writeDB(dbStructure DBStructure) error {
	data, err := json.Marshal(dbStructure)
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = os.WriteFile(db.path, data, 0644)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {

	dbStructure, err := db.loadDB()
	if err != nil && !errors.Is(err, ErrorEmptyFile) {
		log.Fatal(err)
		return Chirp{}, err
	}

	if errors.Is(err, ErrorEmptyFile) {
		dbStructure = DBStructure{
			Chirps: make(map[int]Chirp),
		}
	}

	chirpId := len(dbStructure.Chirps) + 1

	chirp := Chirp{Id: chirpId, Body: body}

	dbStructure.Chirps[chirpId] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		log.Fatal(err)
		return Chirp{}, err
	}

	return chirp, nil

}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil && !errors.Is(err, ErrorEmptyFile) {
		return []Chirp{}, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))

	for _, v := range dbStructure.Chirps {
		chirps = append(chirps, v)
	}

	return chirps, nil

}
