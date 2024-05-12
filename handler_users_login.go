package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rxmeez/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
	}

	user, err := cfg.db.Login(params.Email, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not Authorized")
		return
	}

	defaultExpiration := 60 * 60
	if params.ExpiresInSeconds == 0 {
		params.ExpiresInSeconds = defaultExpiration
	} else if params.ExpiresInSeconds > defaultExpiration {
		params.ExpiresInSeconds = defaultExpiration
	}

	token, err := auth.MakeJWT(user.Id, cfg.jwtSecret, time.Duration(params.ExpiresInSeconds)*time.Second)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create JWT")
		return
	}

	refreshTokenString, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate refresh token")
		return
	}

	err = cfg.db.StoreRefreshToken(user.Id, refreshTokenString, 60*24*time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't store refresh token")
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			Id:          user.Id,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        token,
		RefreshToken: refreshTokenString,
	})

}
