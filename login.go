package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bmkersey/go-chirpy/internal/auth"
	"github.com/bmkersey/go-chirpy/internal/database"
)

func (c *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding params: %s", err)
		sendError(w, 400, "Something went wrong")
		return
	}

	user, err := c.dbQueries.GetUser(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error finding user: %s", err)
		sendError(w, 401, "Incorrect email or password")
		return
	}

	token, err := auth.MakeJWT(user.ID, c.jwtSecret)
	if err != nil {
		log.Printf("Error making JWT token: %s", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Error creating refresh token: %s", err)
		sendError(w, 400, "Something went wrong")
		return
	}
	rToken, err := c.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	})
	if err != nil {
		log.Printf("Error creating refresh token: %s", err)
		return
	}

	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil {
		log.Printf("Passwords do not match: %s", err)
		sendError(w, 401, "Incorrect email or password")
		return
	}

	loggedInUser := User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: rToken.Token,
		IsChirpyRed:  user.IsChirpyRed.Bool,
	}

	sendJsonResponse(w, 200, loggedInUser)
}
