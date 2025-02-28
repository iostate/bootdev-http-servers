package main

import (
	"net/http"

	auth "github.com/iostate/bootdev-http-servers/internal"
)

func (cfg *apiConfig) handleRevoke(w http.ResponseWriter, r *http.Request) {
	authHeader, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error getting Authorization token", err)
		return
	}

	_, err = cfg.db.GetRefreshToken(r.Context(), authHeader)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error retrieving refresh token", err)
		return
	}

	if err := cfg.db.RevokeRefreshToken(r.Context(), authHeader); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error revoking refresh token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
