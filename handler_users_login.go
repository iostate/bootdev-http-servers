package main

import (
	"encoding/json"
	"net/http"

	auth "github.com/iostate/bootdev-http-servers/internal"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	var req UserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding json", err)
		return
	}

	// Get user by email
	user, err := cfg.db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	if err = auth.CheckPasswordHash(req.Password, user.HashedPassword); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	respondWithJSON(w, 200, UserResponse{
		Id: user.ID,
		Created_at: user.CreatedAt,
		Updated_at: user.UpdatedAt,
		Email: user.Email,
	})
}