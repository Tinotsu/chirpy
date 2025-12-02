package main

import (
	"net/http"
	"github.com/google/uuid"
	"github.com/Tinotsu/chirpy/internal/database"
)

type Chirps []Chirp

func (apiCfg * apiConfig) getChirps (w http.ResponseWriter, r *http.Request) {
	chirpsData, err := apiCfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get Chirp from the DB", err)
		return
	}
	chirps := *new(Chirps)
	for _, chirp := range chirpsData {
		chirpStruct := Chirp {
			ID:		chirp.ID,
			CreatedAt:	chirp.CreatedAt,
			UpdatedAt:	chirp.UpdatedAt,
			Body: 	chirp.Body,
			UserID:	chirp.UserID,
		}
		chirps = append(chirps, chirpStruct)
	}
	respondWithJSON(w, http.StatusOK, chirps) 
}

func (apiCfg * apiConfig) getChirpByID (w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")
	chirpsData, err := apiCfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get Chirp from the DB", err)
		return
	}

	isFind := false

	chirp := *new(database.Chirp)
	for _, chirpTarget := range chirpsData {
		if chirpTarget.ID == uuid.MustParse(id) {
			chirp = chirpTarget
			isFind = true
		}
	}

	if !isFind {
		respondWithError(w, http.StatusNotFound, "Chirp ID does not match any ID in the DB", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp {
		ID:		chirp.ID,
		CreatedAt:	chirp.CreatedAt,
		UpdatedAt:	chirp.UpdatedAt,
		Body: 	chirp.Body,
		UserID:	chirp.UserID,
	})
}
