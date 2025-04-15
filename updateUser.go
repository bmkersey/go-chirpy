package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bmkersey/go-chirpy/internal/auth"
	"github.com/bmkersey/go-chirpy/internal/database"
)

func (c *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		sendError(w, 401, "Access token not found")
		return
	}

	accessID, err := auth.ValidateJWT(accessToken, c.jwtSecret)
	if err != nil {
		log.Printf("Error validating user: %s", err)
		sendError(w, 401, "Couldn't validate user")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		sendError(w, 400, "Incorrect or missing parameters")
		return
	}

	hashedPW, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		sendError(w, 400, "Something went wrong")
		return
	}

	updatedUser, err := c.dbQueries.UpdateUser(r.Context(), database.UpdateUserParams{
		HashedPassword: hashedPW,
		Email:          params.Email,
		ID:             accessID,
	})
	if err != nil {
		log.Printf("Error updating user %v: %s", accessID, err)
		sendError(w, 400, "Couldn't update user")
		return
	}

	sendJsonResponse(w, 200, User{
		ID:          updatedUser.ID,
		Email:       updatedUser.Email,
		CreatedAt:   updatedUser.CreatedAt,
		UpdatedAt:   updatedUser.UpdatedAt,
		IsChirpyRed: updatedUser.IsChirpyRed.Bool,
	})
}
