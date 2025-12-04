package main

import (
	"encoding/json"
	"net/http"
	"github.com/Tinotsu/chirpy/internal/auth"
	"github.com/Tinotsu/chirpy/internal/database"
)

func (apiCfg *apiConfig) handlerUpdateUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password  string	`json:"password"`
		Email     string    `json:"email"`
	}
	
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode Chirp", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Couldn't get token from header", err)
		return
	}
	userUUID, err := auth.ValidateJWT(token, apiCfg.envSecret)
	if err != nil {
		respondWithError(w, 401, "Couldn't validate token from user", err)
		return
	}

	userParam := new(database.UpdateUsersParams)
	userParam.Email = params.Email
	userParam.HashedPassword, err = auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password's User", err)
		return
	}
	userParam.ID = userUUID
	err = apiCfg.db.UpdateUsers(r.Context(), *userParam)
	if err != nil {
		respondWithError(w, 401, "Couldn't update email and password from user", err)
		return
	}
	respondWithJSON(w, 200, User {
		Email: params.Email,
	})
}
