package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bmkersey/go-chirpy/internal/auth"
	"github.com/bmkersey/go-chirpy/internal/database"
)


func(c *apiConfig)handlerCreateUser(w http.ResponseWriter, r *http.Request){
	type Paramerters struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := Paramerters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		sendError(w, 400, "Something went wrong")
		return 
	}

	hashedPw, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		sendError(w, 400, "Something went wrong")
		return 
	}

	user, err := c.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hashedPw,
	})
	if err != nil {
		log.Printf("Error creating user: %s", err)
		sendError(w, 400, "Something went wrong while creating user")
		return 
	}

	newUser := User{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
	}

	sendJsonResponse(w, 201, newUser)
}