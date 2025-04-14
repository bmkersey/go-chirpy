package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/bmkersey/go-chirpy/internal/auth"
)


func (c *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
		Password string `json:"password"`
		ExpiresInSeconds *int `json:"expires_in_seconds,omitempty"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding params: %s", err)
		sendError(w, 400, "Something went wrong")
		return
	}
	expiresInSeconds := 3600 // Default 1 hour
	if params.ExpiresInSeconds != nil {
			if *params.ExpiresInSeconds > 0 && *params.ExpiresInSeconds < 3600 {
					expiresInSeconds = *params.ExpiresInSeconds
			}
	}

	

	user, err := c.dbQueries.GetUser(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error finding user: %s", err)
		sendError(w, 401, "Incorrect email or password")
		return
	}

	token, err := auth.MakeJWT(user.ID, c.jwtSecret, time.Duration(expiresInSeconds)*time.Second)
	if err != nil {
		log.Printf("Error making JWT token: %s", err)
		return
	}
	

	err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
	if err != nil{
		log.Printf("Passwords do not match: %s", err)
		sendError(w, 401, "Incorrect email or password")
		return
	}

	loggedInUser := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
	}

	sendJsonResponse(w, 200, loggedInUser)
}