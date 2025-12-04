package main

import (
	"net/http"
	"github.com/google/uuid"
	"github.com/Tinotsu/chirpy/internal/database"
	"sort"
	"strings"
	"log"
)

type Chirps []Chirp

func (apiCfg * apiConfig) getChirps (w http.ResponseWriter, r *http.Request) {
	authorOpt := r.URL.Query().Get("author_id")
	authorOptID,_ := uuid.Parse(authorOpt)

	sortOPT := r.URL.Query().Get("sort")

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
		if authorOptID == chirp.UserID {
			chirps = append(chirps, chirpStruct)
		} else if authorOpt == "" {
			chirps = append(chirps, chirpStruct)
		}
	}

	if strings.Contains(sortOPT, "asc") {
		log.Println("sort is asc")
		sort.Slice(chirps, func (i, j int) bool {return chirps[j].UpdatedAt.After(chirps[i].UpdatedAt)})
	}
	if strings.Contains(sortOPT, "desc") {
		log.Println("sort is desc")
		sort.Slice(chirps, func (i, j int) bool {return chirps[j].UpdatedAt.Before(chirps[i].UpdatedAt)})
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
