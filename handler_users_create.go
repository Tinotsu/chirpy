package main

import (
	"time"
	"github.com/google/uuid"
	"github.com/Tinotsu/chirpy/internal/database"
	"github.com/Tinotsu/chirpy/internal/auth"
	"encoding/json"
	"net/http"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Password  string	`json:"password"`
	Token string		`json:"token"`
	RefreshToken string		`json:"refresh_token"`
	IsChirpyRed bool		`json:"is_chirpy_red"`
}

func (apiCfg *apiConfig) handlerUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password  string	`json:"password"`
		Email     string    `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode User", err)
		return
	}

	userParam := new(database.CreateUserParams)
	userParam.Email = params.Email
	userParam.HashedPassword, err = auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password's User", err)
		return
	}
	user, err := apiCfg.db.CreateUser(r.Context(), *userParam)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create User", err)
		return
	}

	respondWithJSON(w, 201, User {
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.CreatedAt,
		Email: user.Email,
		Password: user.HashedPassword,
		IsChirpyRed: user.IsChirpyRed.Bool,
	})
}
