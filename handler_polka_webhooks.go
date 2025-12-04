package main

import (
	"encoding/json"
	"net/http"
	"github.com/google/uuid"
	"github.com/Tinotsu/chirpy/internal/auth"
)

func (apiCfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	polkaKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "Couldn't get polkaKey from header", err)
		return
	}
	if polkaKey != apiCfg.polkaAPI {
		respondWithError(w, 401, "Couldn't verify polkaKey from header", err)
		return
	}

	type userParams struct {
		UserID	uuid.UUID	`json:"user_id"`
	}
	type parameters struct {
		Event  string	`json:"event"`
		Data 	userParams	
	}    

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode polka webhooks", err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, 204, "")
	} else {
		err = apiCfg.db.UpgradeChirpyRedByID(r.Context(), params.Data.UserID)
		if err != nil {
			respondWithError(w, 404, "Couldn't find User with this ID", err)
			return
		}
		respondWithJSON(w, 204, "")
	}
}
