package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bmkersey/go-chirpy/internal/auth"
	"github.com/google/uuid"
)

func (c *apiConfig) handlerUpgradeUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		log.Printf("Error getting API key from headers: %s", err)
		sendError(w, 401, "Invalid or malformed API key")
		return
	}

	if apiKey != c.polkaKey {
		log.Println("Unauthorized attempt to upgrade user")
		sendError(w, 401, "You are not authorized to use this resource")
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding params: %s", err)
		sendError(w, 400, "Something went wrong")
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(204)
		return
	}

	userID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		log.Printf("Error converting user id string to UUID: %s", err)
		sendError(w, 400, "Invalid UserID")
		return
	}

	err = c.dbQueries.UpgradeUser(r.Context(), userID)
	if err != nil {
		log.Printf("Error upgrading user: %s", err)
		sendError(w, 404, "Could not find user")
		return
	}

	w.WriteHeader(204)
}
