package main

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {

	authorID := r.URL.Query().Get("author_id")
	if authorID != "" {
		// filter chirps by author_id
		authorID, err := uuid.Parse(authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Can't parse authorID into UUID", err)
			return
		}
		chirpsByAuthor, err := cfg.db.GetChirpsByUserId(r.Context(), authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting chirps By Author ID %s", authorID), err)
			return
		}

		var chirps []Chirp

		for _, chirp := range chirpsByAuthor {
			addChirp := Chirp{
				Id:        chirp.ID,
				CreatedAt: chirp.CreatedAt,
				UpdatedAt: chirp.UpdatedAt,
				Body:      chirp.Body,
				UserID:    chirp.UserID,
			}
			chirps = append(chirps, addChirp)
		}

		// Build response
		respondWithJSON(w, http.StatusOK, chirps)
		return
	}
	var chirpsResponse []Chirp

	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving all chirps: ", err)
		return
	}

	for _, chirp := range chirps {
		chirpsResponse = append(chirpsResponse, Chirp{
			Id:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJSON(w, 200, chirpsResponse)
	return

}
