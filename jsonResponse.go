package main

import (
	"encoding/json"
	"log"
	"net/http"
)


func sendJsonResponse(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		sendError(w, 400, "Something went wrong")
		return 
	}
	w.Write(data)
}