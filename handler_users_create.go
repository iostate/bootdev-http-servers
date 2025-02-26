package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	auth "github.com/iostate/bootdev-http-servers/internal"
	"github.com/iostate/bootdev-http-servers/internal/database"
)

type UserRequest struct {
	Email 		string 		`json:"email"`
	Password	string		`json:"password"`
}

type UserResponse struct {
	Id 			uuid.UUID	`json:"id"`
	Created_at 	time.Time 	`json:"created_at"`
	Updated_at	time.Time 	`json:"updated_at"`
	Email 		string 		`json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	var req UserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("Error decoding json: %v", err)
		respondWithError(w, 500, "Error decoding json: ", err)
		return
	}

	// Hash password
	hashedPw, err := auth.HashPassword(req.Password)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}
	
	// Create User
	
	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email: req.Email, 
		HashedPassword: hashedPw,
	})
	if err != nil {
		log.Printf("Error creating user: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error creating user", err)
		return
	}

	// Build JSON response
	userResponse := UserResponse{
		Id: user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email: user.Email,
	}
	respondWithJSON(w,http.StatusCreated, userResponse)
}