package main

import (
	"log"
	"net/http"
	"time"

	auth "github.com/iostate/bootdev-http-servers/internal"
)

// Should return a new refresh token for that user in response
func (cfg *apiConfig) handleRefresh(w http.ResponseWriter, r *http.Request) {
	type tokenResponse struct {
		AccessToken string `json:"token"`
	}
	authHeader, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 400, "Error getting token", err)
		return
	}

	// If token doesn't exist, return 401
	getRefreshToken, err := cfg.db.GetRefreshToken(r.Context(), authHeader)
	if err != nil {
		errorFindingTokenStr := "error finding refresh token"
		log.Printf("%s: %s", errorFindingTokenStr, authHeader)
		respondWithError(w, http.StatusUnauthorized, errorFindingTokenStr, err)
		return
	}

	// If refresh token is expired, return 401
	if getRefreshToken.ExpiresAt.Before(time.Now()) {
		tokenExpiredSt := "refresh token expired"
		log.Printf("%s: %s", tokenExpiredSt, authHeader)
		respondWithError(w, http.StatusUnauthorized, tokenExpiredSt, nil)
		return
	}

	// Checks to see if RevokedAt is populated,
	// if it is, then this refresh token is trash and we should return 401
	if getRefreshToken.RevokedAt.Valid {
		tokenRevokedStr := "refresh token has been revoked"
		log.Printf("%s: %s", tokenRevokedStr, authHeader)
		respondWithError(w, http.StatusUnauthorized, tokenRevokedStr, nil)
		return
	}

	// If Couldn't find refresh token by user, return 401
	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), authHeader)
	// if token doesn't exist, respond with 401
	if err != nil {
		log.Printf("Couldnt find refresh token using user id of %s", authHeader)
		respondWithError(w, 401, "Could not find refresh token", err)
		return
	}

	// Create new access token
	newAccessToken, err := auth.MakeJWT(user.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		log.Printf("Error making new access token")
		respondWithError(w, 500, "Error making new access token", err)
		return
	}
	// respond with 200 with new access token
	respondWithJSON(w, 200, tokenResponse{
		AccessToken: newAccessToken,
	})

}
