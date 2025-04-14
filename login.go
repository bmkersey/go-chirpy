package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bmkersey/go-chirpy/internal/auth"
)


func (c *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
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
	}

	sendJsonResponse(w, 200, loggedInUser)
}