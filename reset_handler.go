package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	// if platform != dev, return 403
	if cfg.platform != "dev" {
		err := fmt.Errorf("Currently in dev mode, certain actions are prohibited")
		respondWithError(w, http.StatusForbidden, "FORBIDDEN - Can't delete users while in dev mode", err)
	}
		// reset to zero
		cfg.fileServerHits.Store(0)
		cfg.db.DeleteUsers(r.Context())
		w.Write([]byte("Counter reset to 0 and database reset to initial state"))
}