package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/rxmeez/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerChirpDeleteId(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respondWithError(w, 500, "Id is not a int")
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT")
		return
	}

	userId, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT")
		return
	}

	authorId, err := strconv.Atoi(userId)
	if err != nil {
		log.Fatal("Failed to convert string to int authorId")
	}

	err = cfg.db.DeleteChirp(id, authorId)
	if err != nil {
		respondWithError(w, http.StatusForbidden, "Forbidden to delete")
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)

}
