package main

import (
	"log"
	"net/http"
)


func (c *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := c.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("Error retrieving chirps: %s", err)
		sendError(w, 400, "Something went wrong")
	}

	chirpsRes := []Chirp{}
	for _, chirp := range chirps{
		chirpsRes = append(chirpsRes,Chirp{
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			User_id: chirp.UserID,

		})
	}

	sendJsonResponse(w, 200, chirpsRes)
}