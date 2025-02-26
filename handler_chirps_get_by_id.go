package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirpsById(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Printf("Error parsing UUID: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error parsing UUID", err)
		return
	}
	chirp, err := cfg.db.GetChirpById(r.Context(), chirpID)
	if err != nil {
		log.Printf("Error getting chirp by ID: %v", err)
		respondWithError(w, http.StatusNotFound, "Error retrieving chirp by ID", err)
		return
	}

	respondWithJSON(w, 200, Chirp{
		Id: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})






}