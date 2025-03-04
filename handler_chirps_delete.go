package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request) {

	chirpIDString := r.PathValue("chirpID")
	chirpUUID, err := uuid.Parse(chirpIDString)

	// get user id
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Failed to get authUserID from context"))
		return
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing chirpID into UUID", err)
		return
	}
	chirp, err := cfg.db.GetChirpById(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Chirp not found", err)
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "User is not the owner of the chirp", err)
		return
	}

	// we've checked that the chirp owner is the same as the one
	// from the access token, now it is time to delete the chirp
	if err := cfg.db.DeleteChirpById(r.Context(), chirpUUID); err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error deleting the chirp with ID of %v\n", chirpUUID), err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
