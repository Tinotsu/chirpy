package main

import (
	"errors"
	"fmt"
	"time"
	"github.com/Tinotsu/chirpy/internal/auth"
	"encoding/json"
	"net/http"
)

func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password  string	`json:"password"`
		Email     string    `json:"email"`
		ExpiresInSecond	int	`json:"expires_in_second"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode User", err)
		return
	}

	user, err := apiCfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password's User", err)
		return
	}
	
	isAuth, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't check password's User", err)
		return
	}

	if !isAuth {
		err = errors.New("incorrect email or password")
		respondWithError(w, 401, "Incorrect email or password", err)
		return
	}

	var expTimeSec int

	if params.ExpiresInSecond == 0{
		expTimeSec = 3600
	} else if  params.ExpiresInSecond > 3600 {
		expTimeSec = 3600
	} else {
		expTimeSec = params.ExpiresInSecond
	}
	timeinstring := fmt.Sprint(expTimeSec)

	timeintime, err := time.ParseDuration(fmt.Sprintf("%ss",timeinstring))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't convert timeinstring to timeintime", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, apiCfg.envSecret, timeintime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get the token", err)
		return
	}

	respondWithJSON(w, 200, User {
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.CreatedAt,
		Email: user.Email,
		Token: token,
	})
}
