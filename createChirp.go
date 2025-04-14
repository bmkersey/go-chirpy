package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"unicode"

	"github.com/bmkersey/go-chirpy/internal/auth"
	"github.com/bmkersey/go-chirpy/internal/database"
)


func (c *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding params: %s", err)
		sendError(w, 400, "Something went wrong")
		return 
	}

	token,err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting Bearer token: %s", err)
		sendError(w, 401, "Could not verify user")
		return
	}
	tokenID, err := auth.ValidateJWT(token, c.jwtSecret)
	if err != nil {
		log.Printf("Error validating token: %s", err)
		sendError(w, 401, "Could not verify user")
		return
	}

	if len(params.Body) > 140 {
		log.Println("Chirp is too long")
		sendError(w, 400, "Chirp is too long")
		return
	}

	params.Body = removeProfanity(params.Body)

	chirp, err := c.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   params.Body,
		UserID: tokenID,
	})
	if err != nil {
		log.Printf("Error creating chirp: %s", err)
		sendError(w, 400, "Something went wrong while creating chirp")
		return
	}

	newChirp := Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		User_id: chirp.UserID,
	}

	sendJsonResponse(w, 201, newChirp)

}

func removeProfanity(unclean string) string {
	split := strings.Split(unclean, " ")
	forbidden := map[string]bool{"kerfuffle": true, "sharbert": true, "fornax": true}
	for i,word := range split {
		var cleaned strings.Builder
    for _, char := range word {
      if unicode.IsLetter(char) || unicode.IsDigit(char) {
        cleaned.WriteRune(char)
      }
    }
		if _, ok := forbidden[strings.ToLower(cleaned.String())]; ok {
			split[i] = "****"
		}
	}
	
	return strings.Join(split, " ")
}