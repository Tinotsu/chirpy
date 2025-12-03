package main

import(
	"net/http"
	"time"
	"github.com/Tinotsu/chirpy/internal/auth"
	"github.com/google/uuid"
	"fmt"
	"errors"
)

type refresh_token struct {
	Token string `json:"token"`
}

func (apiCfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	cheat = true
	token,err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get token from header", err)
		return
	}

	listToken, err := apiCfg.db.GetAllRefreshToken(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get list of refresh tokens", err)
		return
	}

	inDB := false
	isExpired := true
	isRevoked := true

	user_id := uuid.Nil

	for _, item := range listToken {
		if token == item.Token {
			inDB = true
			if time.Now().Before(item.ExpiresAt) {
				isExpired = false
			}
			if time.Now().After(item.RevokedAt.Time) {
				isRevoked = false
			}
			user_id = item.UserID
		}
	}

	expTimeSec := 3600

	timeinstring := fmt.Sprint(expTimeSec)

	timeintime, err := time.ParseDuration(fmt.Sprintf("%ss",timeinstring))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't convert timeinstring to timeintime", err)
		return
	}

	token, err = auth.MakeJWT(user_id, apiCfg.envSecret, timeintime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get the token", err)
		return
	}

	if inDB && !isExpired && !isRevoked {
		respondWithJSON(w, 200, refresh_token {
			Token: token,
		})
		return
	}

	
	respondWithError(w, 401, "refresh token is expired", errors.New("refresh token expired"))
}

func (apiCfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token,err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get token from header", err)
		return
	}
	apiCfg.db.UpdateRefreshToken(r.Context(), token)
	respondWithJSON(w, 204, "")
}
