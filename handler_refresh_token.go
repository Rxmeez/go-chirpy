package main

import (
	"net/http"
	"time"

	"github.com/rxmeez/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find refresh token")
	}

	type response struct {
		JWTToken string `json:"token"`
	}

	userId, err := cfg.db.ValidateRefreshToken(refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate refresh token")
		return
	}

	// refreshToken, err = auth.MakeRefreshToken()
	// if err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, "Couldn't create Refresh Token")
	// 	return
	// }

	defaultExpiration := 60 * 60

	token, err := auth.MakeJWT(userId, cfg.jwtSecret, time.Duration(defaultExpiration)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		JWTToken: token,
	})
}
