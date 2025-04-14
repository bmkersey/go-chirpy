package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorResponse struct {
  Error string `json:"error"`
}


func sendError(w http.ResponseWriter, status int, errorMsg string) {
  respBody := errorResponse{
    Error: errorMsg,
  }

  data, err := json.Marshal(respBody)
  if err != nil {
    log.Printf("Error marshalling error response: %s", err)
    w.WriteHeader(500)
    w.Write([]byte(`{"error":"Internal server error"}`))
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(status)
  w.Write(data)
}