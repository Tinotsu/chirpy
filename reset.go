package main

import (
	"net/http"
	"context"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.env == "dev" {
		ctx := context.Background()
		cfg.fileserverHits.Store(0)
		cfg.db.DeleteUsersTable(ctx)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hits reset to 0\nusers table reset"))
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}
