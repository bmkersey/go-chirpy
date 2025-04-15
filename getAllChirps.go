package main

import (
	"log"
	"net/http"
	"sort"

	"github.com/bmkersey/go-chirpy/internal/database"
	"github.com/google/uuid"
)

func (c *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")
	sortOrder := r.URL.Query().Get("sort")
	var chirps []database.Chirp
	var err error

	if authorID == "" {
		chirps, err = c.dbQueries.GetAllChirps(r.Context())
		if err != nil {
			log.Printf("Error retrieving chirps: %s", err)
			sendError(w, 400, "Something went wrong")
		}
	} else {
		authorUUID, err := uuid.Parse(authorID)
		if err != nil {
			sendError(w, 400, "Could not turn provided author ID to UUID")
			return
		}
		chirps, err = c.dbQueries.GetChirpsByAuthor(r.Context(), authorUUID)
		if err != nil {
			log.Printf("Error retrieving chirps: %s", err)
			sendError(w, 400, "Something went wrong")
		}
	}

	chirpsRes := []Chirp{}
	for _, chirp := range chirps {
		chirpsRes = append(chirpsRes, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			User_id:   chirp.UserID,
		})
	}

	if sortOrder == "desc" {
		sort.Slice(chirpsRes, func(i, j int) bool {
			return chirpsRes[j].CreatedAt.Before(chirpsRes[i].CreatedAt)
		})
	}

	sendJsonResponse(w, 200, chirpsRes)
}
