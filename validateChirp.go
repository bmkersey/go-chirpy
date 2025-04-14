package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"unicode"
)


func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		sendError(w, 400, "Something went wrong")
		return 
	}

	if len(params.Body) > 140 {
		log.Println("Chirp is too long")
		sendError(w, 400, "Chirp is too long")
		return
	}

	params.Body = removeProfanity(params.Body)

	type returnVals struct{
		Cleaned_Body string `json:"cleaned_body"`
	}

	respBody := returnVals{
		Cleaned_Body: params.Body,
	}

	

	sendJsonResponse(w, 200, respBody)
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