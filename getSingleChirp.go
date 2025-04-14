package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)


func (c *apiConfig) handlerGetSingleChirp(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Printf("Could not parse id to UUID: %s", err)
		sendError(w, 400, "Something went wrong")
		return
	}

	chirp, err := c.dbQueries.GetSingleChirp(r.Context(), id)
	if err != nil {
		log.Printf("Could not find chirp: %s", err)
		sendError(w, 404, "Chould not find chirp")
		return
	}

	foundChirp := Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		User_id: chirp.UserID,
	}

	sendJsonResponse(w, 200, foundChirp)
}