package main

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {

	// QUERY BY AUTHOR ID
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
	// END QUERY BY AUTHOR ID

	var chirpsResponse []Chirp
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error retrieving all chirps: ", err)
		return
	}

	// QUERY SORT BY ASC OR DESC
	sortParam := r.URL.Query().Get("sort")
	if sortParam != "" {
		if strings.ToLower(sortParam) == "asc" {
			sort.Slice(chirps, func(i, j int) bool {
				return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
			})
		}

		if strings.ToLower(sortParam) == "desc" {
			sort.Slice(chirps, func(i, j int) bool {
				return chirps[j].CreatedAt.Before(chirps[i].CreatedAt)
			})
		}
	}
	// END QUERY SORT BY ASC OR DESC

	// CONVERT TO CHIRP STRUCT WITH JSON TAGS
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
