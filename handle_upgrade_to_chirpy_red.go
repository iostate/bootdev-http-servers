package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

/*
		Accepts req.body of form:
		{
			"event": "user.upgraded",
			"data": {
				"user_id": "3311741c-680c-4546-99f3-fc9efac2036c"
	  }
	}
*/
func (cfg *apiConfig) handleUpgradeToChirpyRed(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	params := &parameters{}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding JSON", err)
		return
	}

	// if event is anything but "user.upgraded", return 204 and end
	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Attempt to find user, if can't find user, return 404
	_, err := cfg.db.GetUserById(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Error finding user in database with ID of %s", params.Data.UserID), err)
		return
	}

	// means we found user, now let's make ANOTHER call to DB and update the user to Chirpy Red
	userRecord, err := cfg.db.UpdateUserToChirpyRed(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error upgrading user to Chirpy Red", err)
		return
	}

	// if user is successfully upgraded, now we just return a 204
	w.WriteHeader(204)
}
