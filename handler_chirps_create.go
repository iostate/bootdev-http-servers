package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/iostate/bootdev-http-servers/internal/database"
)

type requestBody struct {
	Body 	string `json:"body"`
	UserID 	uuid.UUID 	`json:"user_id"`
}

// Chirp response
type Chirp struct {
	Id 			uuid.UUID 	`json:"id"`
	CreatedAt 	time.Time 	`json:"created_at"`
	UpdatedAt 	time.Time	`json:"updated_at"`
	Body		string		`json:"body"`
	UserID		uuid.UUID	`json:"user_id"`
}


func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {

	var reqBody requestBody
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		log.Printf("Error decoding json: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error decoding json: ", err)
		return
	}
	if len(reqBody.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is longer than 140 characters", err)
		return
	}
	
	// TODO:
	// Start checking for valid UUID? Let's start without checks, 
	// since it'll be a string and I don't know how the tests are written.
	// I am getting they are using generate random UUID, though.


	// Sanitize bad words - if we get to here, everything is correct,
	cleanedChirp := replaceBadWords(reqBody.Body)
	// let's create the chirp
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleanedChirp,
		UserID: reqBody.UserID,
	})
	if err != nil {
		log.Printf("Error creating the chirp in DB: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error creating the chirp in DB", err)
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		Body: chirp.Body,
		Id: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		UserID: chirp.UserID,
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
		if unicode.IsPunct(rune(word[len(word) - 1])) {
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