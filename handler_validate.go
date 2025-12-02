package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerChirpsValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		Valid bool `json:"valid"`
		CleanBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	listBody := strings.Split(params.Body," ")
	for i, w := range listBody {
		p := strings.ToLower(w)
		if p == "kerfuffle" || p == "sharbert" || p == "fornax" {
			listBody[i] = "****"
		}
	}
	newBody := strings.Join(listBody, " ")

	respondWithJSON(w, http.StatusOK, returnVals{
		Valid: true,
		CleanBody: newBody,
	})
}
