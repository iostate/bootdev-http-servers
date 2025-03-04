package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	auth "github.com/iostate/bootdev-http-servers/internal"
	"github.com/iostate/bootdev-http-servers/internal/database"
)

type User struct {
	ID             uuid.UUID `json:"id`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	Email          string
	HashedPassword string
}

type UserUpdateResponse struct {
	ID         uuid.UUID `json:"id"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
	Email      string    `json:"email"`
	ChirpyRed  bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	// storeUserInContextMiddleware sets the user ID in context
	authUserID, ok := r.Context().Value("userID").(uuid.UUID)

	if !ok {
		w.WriteHeader(401)
		w.Write([]byte("Failed to get authUserID from context"))
		return
	}

	params := &parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, 500, "error decoding email and password", err)
		return
	}

	newHashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "error hashing password", err)
		return
	}

	userRecord, err := cfg.db.UpdateUserPasswordAndEmail(r.Context(), database.UpdateUserPasswordAndEmailParams{
		ID:             authUserID,
		Email:          params.Email,
		HashedPassword: newHashedPassword,
	})

	type response struct {
		UserID uuid.UUID `json:"user_id"`
	}

	respondWithJSON(w, 200, UserUpdateResponse{
		ID:         userRecord.ID,
		Email:      userRecord.Email,
		Created_at: userRecord.CreatedAt,
		Updated_at: userRecord.UpdatedAt,
	})
}
