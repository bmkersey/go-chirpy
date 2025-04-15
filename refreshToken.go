package main

import (
	"log"
	"net/http"

	"github.com/bmkersey/go-chirpy/internal/auth"
)

func (c *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Couldn't find token")
		return
	}

	user, err := c.dbQueries.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		log.Printf("Error retrieving user from refresh token: %s", err)
		sendError(w, 401, "Could not find user with provided token")
		return
	}

	newAccessToken, err := auth.MakeJWT(user.ID, c.jwtSecret)
	if err != nil {
		log.Printf("Error creating access token: %s", err)
		sendError(w, 401, "Could not create new access token")
		return
	}

	sendJsonResponse(w, 200, response{
		Token: newAccessToken,
	})
}
