package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/rxmeez/chirpy/internal/database"
)

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {

	db := cfg.db

	if r.Method == "GET" {
		chirps, err := db.GetChirps()
		if err != nil {
			if !errors.Is(err, database.ErrorEmptyFile) {
				log.Printf("Error decoding parameters: %s", err)
				w.WriteHeader(500)
				return
			}
		}

		fmt.Println(chirps)

		chirpBytes, err := json.Marshal(chirps)
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(chirpBytes)
		return
	}

	type parameters struct {
		Body string `json:"body"`
	}

	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
	}

	if len(params.Body) > 140 {
		dat, err := json.Marshal(struct {
			Error string `json:"error"`
		}{
			Error: "Chirp is too long",
		})
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	words := strings.Split(params.Body, " ")
	var cleaned_body []string

w:
	for _, word := range words {
		for _, profaneWord := range profaneWords {
			if strings.ToLower(word) == profaneWord {
				cleaned_body = append(cleaned_body, "****")
				continue w
			}
		}
		cleaned_body = append(cleaned_body, word)
	}

	body := strings.Join(cleaned_body, " ")

	chirp, err := db.CreateChirp(body)
	if err != nil {
		log.Printf("Error Creating Chirp: %s", err)
		w.WriteHeader(500)
		return
	}

	chirpBytes, err := json.Marshal(chirp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(201)
	w.Write(chirpBytes)
	return

}
