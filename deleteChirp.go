package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/bmkersey/go-chirpy/internal/auth"
	"github.com/google/uuid"
)

func (c *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Printf("Error getting chirpID from url: %s", err)
		sendError(w, 400, "Incorrect chirpID")
		return
	}

	chirp, err := c.dbQueries.GetSingleChirp(r.Context(), chirpID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			sendError(w, 404, "Chirp not found")
		} else {
			log.Printf("Error getting chirp: %s", err)
			sendError(w, 500, "Internal server error")
		}
		return
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error decoding access token: %s", err)
		sendError(w, 401, "Access token must be in format: Bearer <token>")
		return
	}

	accessID, err := auth.ValidateJWT(accessToken, c.jwtSecret)
	if err != nil {
		log.Printf("Error validating JWT: %s", err)
		sendError(w, 401, "Invalid JWT")
		return
	}

	if accessID != chirp.UserID {
		log.Printf("Unauthorized attemp to delete chirp")
		sendError(w, 403, "You do not have permisson to delete this chirp")
		return
	}

	err = c.dbQueries.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		log.Printf("Error deleting chirp: %s", err)
		sendError(w, 403, "You are not the owner of this chirp")
		return
	}

	w.WriteHeader(204)
}
