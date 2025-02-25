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

	type cleanedBody struct {
		CleanedBody string `json:"cleaned_body"`
	}

	var reqBody requestBody
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		log.Printf("error decoding json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ := json.Marshal(errorResponse{Error: "Something went wrong"})
		w.Write(resp)
		return
	}

	if len(reqBody.Body) > 140 {
		w.WriteHeader(http.StatusBadRequest)
		resp, err := json.Marshal(errorResponse{Error: "Chirp is too long"})
		if err != nil {
			log.Printf("Error marshalling error response: %v", err)
			w.Write([]byte(`"error":"Something went wrong`))
			return
		}
		w.Write(resp)
		return
	} 

	// Sanitize bad words
	cleanedChirp := replaceBadWords(reqBody.Body)

	// Chirp is valid
	// resp, err := json.Marshal(validResponse{Valid: true})
	resp, err := json.Marshal(cleanedBody{CleanedBody: cleanedChirp})
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`"error":"Something went wrong`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
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