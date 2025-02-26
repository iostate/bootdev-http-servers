package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type EmailRequest struct {
	Email 	string 		`json:"email"`
}

type UserResponse struct {
	Id 	uuid.UUID	`json:"id"`
	Created_at 	time.Time `json:"created_at"`
	Updated_at	time.Time `json:"updated_at"`
	Email 	string `json:"email"`

}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	var req EmailRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Error decoding json: %v", err)
		respondWithError(w, 500, "Error decoding json: ", err)
		return
	}

	log.Printf("Email: %v\n", req.Email)

	// Create user
	user, err := cfg.db.CreateUser(r.Context(), req.Email)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		respondWithError(w, 500, "Error creating user: ", err)
		return
	}

	userResponse := UserResponse{
		Id: user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email: user.Email,
	}

	respondWithJSON(w,http.StatusCreated, userResponse)
}