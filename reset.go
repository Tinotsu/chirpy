package main

import (
	"net/http"
	"context"
)

func (apiCfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if apiCfg.env == "dev" {
		ctx := context.Background()
		apiCfg.fileserverHits.Store(0)
		apiCfg.db.DeleteUsersTable(ctx)
		apiCfg.db.DeleteChirpsTable(ctx)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hits reset to 0\nusers table reset"))
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}
