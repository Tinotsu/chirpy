package main

import (
	"net/http"
	"github.com/Tinotsu/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerDeleteChirps (w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Couldn't get token from header", err)
		return
	}
	userID, err := auth.ValidateJWT(token, apiCfg.envSecret)
	if err != nil {
		respondWithError(w, 403, "Couldn't delete, user is not the author", err)
		return
	}

	chirpsData, err := apiCfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get Chirp from the DB", err)
		return
	}

	isAllowed := false

	for _, chirpTarget := range chirpsData {
		if chirpTarget.ID == uuid.MustParse(id) && chirpTarget.UserID == userID {
			isAllowed = true
		}
	}

	if !isAllowed {
		respondWithError(w, 403, "Couldn't delete, user is not the author", err)
		return
	}

	err = apiCfg.db.DeleteChirp(r.Context(), uuid.MustParse(id))
	if err != nil {
		respondWithError(w, 404, "Couldn't delete, chirp not found", err)
		return
	}

	respondWithJSON(w, 204, "") 
}
