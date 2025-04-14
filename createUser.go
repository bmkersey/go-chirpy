package main

import (
	"encoding/json"
	"log"
	"net/http"
)


func(c *apiConfig)handlerCreateUser(w http.ResponseWriter, r *http.Request){
	type Paramerters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := Paramerters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		sendError(w, 400, "Something went wrong")
		return 
	}

	user, err := c.dbQueries.CreateUser(r.Context(), params.Email)
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