package main

import (
	"encoding/json"
	"net/http"

	"github.com/rxmeez/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerUsersUpgrade(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Data struct {
			UserId int `json:"user_id"`
		} `json:"data"`
		Event string `json:"event"`
	}

	polka, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find Polka ApiKey")
		return
	}

	if polka != cfg.polkaSecret {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized Polka apikey")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, nil)
		return
	}

	user, err := cfg.db.UpgradeUser(params.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		Id:          user.Id,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})

}
