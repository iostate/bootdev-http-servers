package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"unicode"
)

func ValidateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Body 	string `json:"body"`
	}

	type errorResponse struct{
		Error 	string `json:"error"`
	}

	type validResponse struct {
		Valid bool `json:"valid"`
	}

	type cleanedResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

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

	// Sanitize bad words
	cleanedChirp := replaceBadWords(reqBody.Body)
	respondWithJSON(w, http.StatusOK, cleanedResponse{CleanedBody: cleanedChirp})
}

// Replaces any bad words, words are preset in the function
// "kerfuffle", "sharbert", "fornax"
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