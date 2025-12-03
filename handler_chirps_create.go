package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"fmt"
	"github.com/google/uuid"
	"github.com/Tinotsu/chirpy/internal/database"
	"github.com/Tinotsu/chirpy/internal/auth"
	"time"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body     string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (apiCfg *apiConfig) handlerChirps (w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
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
		respondWithError(w, http.StatusInternalServerError, "Couldn't get token from header", err)
		return
	}
	userUUID, err := auth.ValidateJWT(token, apiCfg.envSecret)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't validate token from user", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cleaned := getCleanedBody(params.Body, badWords)

	chirpParam := new(database.CreateChirpParams)
	chirpParam.Body = cleaned
	chirpParam.UserID = userUUID
	chirp, err := apiCfg.db.CreateChirp(r.Context(), *chirpParam)
	fmt.Println("\n\nUSER ID : ", chirp.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create Chirp", err)
		return
	}

	respondWithJSON(w, 201, Chirp {
		ID:		chirp.ID,
		CreatedAt:	chirp.CreatedAt,
		UpdatedAt:	chirp.UpdatedAt,
		Body: 	chirp.Body,
		UserID:	chirp.UserID,
	})
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
