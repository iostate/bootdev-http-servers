package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	auth "github.com/iostate/bootdev-http-servers/internal"
	"github.com/iostate/bootdev-http-servers/internal/database"
)

// Chirp response
type Chirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	// AUTH
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error getting bearer token", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error validating JWT", err)
		return
	}

	// Now that auth is done, then decode
	params := &parameters{}
	err = json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		log.Printf("Error decoding json: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error decoding json: ", err)
		return
	}

	// Chirp is too long
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is longer than 140 characters", err)
		return
	}

	// Sanitize bad words - if we get to here, everything is correct,
	cleanedChirp := replaceBadWords(params.Body)
	// let's create the chirp
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedChirp,
		UserID: userID,
	})
	if err != nil {
		log.Printf("Error creating the chirp in DB: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error creating the chirp in DB", err)
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		Body:      chirp.Body,
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		UserID:    chirp.UserID,
	})

}

// Replaces any bad words, words are preset in the function
// "kerfuffle", "sharbert", "fornax"
// P.S. my function naming is awesome
func replaceBadWords(chirpMsg string) (cleanedString string) {
	unacceptableWords := []string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Split(chirpMsg, " ")

	// Loop through string to check for bad words
	for i, word := range words {

		// Check for empty words
		if len(word) == 0 {
			continue
		}

		// Check for punctuation
		if unicode.IsPunct(rune(word[len(word)-1])) {
			continue
		}

		// Check to see if word is an unacceptable word
		for _, unacceptableWord := range unacceptableWords {
			if strings.Contains(strings.ToLower(word), unacceptableWord) {
				words[i] = "****"
				break
			}
		}

	}

	return strings.Join(words, " ")
}
