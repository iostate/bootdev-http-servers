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

type JWTUserRequest struct {
	Email            string        `json:"email"`
	Password         string        `json:"password"`
	ExpiresInSeconds time.Duration `string:"expires_in_seconds"`
}

type UserLoginResponse struct {
	Id           uuid.UUID `json:"id"`
	Created_at   time.Time `json:"created_at"`
	Updated_at   time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	JWTToken     string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		// ExpiresInSeconds time.Duration `string:"expires_in_seconds"`
	}
	params := &parameters{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding json", err)
		return
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
	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour) // explicitly set to 1 hour
	if err != nil {
		respondWithError(w, 500, "Error creating JWT", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Error creating refresh token")
		respondWithError(w, 500, "Error creating refresh token", err)
		return
	}

	refreshTokenDbItem, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, 500, "Trouble creating refresh token in db", err)
		return
	}

	respondWithJSON(w, 200, UserLoginResponse{
		Id:           user.ID,
		Created_at:   user.CreatedAt,
		Updated_at:   user.UpdatedAt,
		Email:        user.Email,
		JWTToken:     token,
		RefreshToken: refreshTokenDbItem.Token,
		IsChirpyRed:  user.IsChirpyRed.Bool,
	})
}
