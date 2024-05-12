package main

import (
	"errors"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/rxmeez/chirpy/internal/database"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {

	dbChirps, err := cfg.db.GetChirps()
	if err != nil && !errors.Is(err, database.ErrorEmptyFile) {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	chirps := []Chirp{}

	authorId := r.URL.Query().Get("author_id")
	if authorId != "" {
		for _, dbChirp := range dbChirps {
			authorIdInt, err := strconv.Atoi(authorId)
			if err != nil {
				log.Fatal("Unable to convert authorId string to int")
				return
			}
			if authorIdInt == dbChirp.AuthorId {
				chirps = append(chirps, Chirp{
					Id:       dbChirp.Id,
					Body:     dbChirp.Body,
					AuthorId: dbChirp.AuthorId,
				})
			}
		}
	} else {
		for _, dbChirp := range dbChirps {
			chirps = append(chirps, Chirp{
				Id:       dbChirp.Id,
				Body:     dbChirp.Body,
				AuthorId: dbChirp.AuthorId,
			})
		}
	}

	sorter := "asc"
	if r.URL.Query().Get("sort") == "desc" {
		sorter = "desc"
	}

	sort.Slice(chirps, func(i, j int) bool {
		if sorter == "asc" {
			return chirps[i].Id < chirps[j].Id
		}
		return chirps[i].Id > chirps[j].Id
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerChirpRetrieveId(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondWithError(w, 500, "Id is not a int")
	}

	chirp, err := cfg.db.GetChirp(id)
	if err != nil && !errors.Is(err, database.ErrorEmptyFile) {
		respondWithError(w, http.StatusNotFound, "Couldn't retrieve chirps")
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)

}
