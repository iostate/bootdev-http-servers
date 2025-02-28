package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	auth "github.com/iostate/bootdev-http-servers/internal"
)

// type JWTUserRequest struct {
// 	UserRequest
// 	ExpiresInSeconds	int	`string:"expires_in_seconds"`
// }

type JWTUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	// Might need to be int
	ExpiresInSeconds time.Duration `string:"expires_in_seconds"`
}

type UserLoginResponse struct {
	Id         uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Email      string    `json:"email"`
	JWTToken   string    `json:"token"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string        `json:"email"`
		Password         string        `json:"password"`
		ExpiresInSeconds time.Duration `string:"expires_in_seconds"`
	}
	params := &parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding json", err)
		return
	}

	if params.ExpiresInSeconds == 0 || params.ExpiresInSeconds > time.Hour {
		params.ExpiresInSeconds = time.Hour // should this be int?
	}

	// Get user by email
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	if err = auth.CheckPasswordHash(params.Password, user.HashedPassword); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// Generate JWT and give it back
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, params.ExpiresInSeconds)
	if err != nil {
		respondWithError(w, 500, "Error creating JWT", err)
		return
	}
	// Store cookie in response
	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization", // or a name of your choice
		Value:    token,
		Expires:  time.Now().UTC().Add(params.ExpiresInSeconds),
		HttpOnly: true,  // enhances security by not allowing JavaScript access
		Secure:   false, // if you're using HTTPS
		Path:     "/",   // ensure it matches your application's path
	})

	respondWithJSON(w, 200, UserLoginResponse{
		Id:         user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email:      user.Email,
		JWTToken:   token,
	})
}
