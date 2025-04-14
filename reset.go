package main

import (
	"log"
	"net/http"
)

func (c *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	c.fileserverHits.Store(0)
	if c.platform != "dev"{
		sendError(w, 403, "Forbidden")
		return
	}
	err := c.dbQueries.DeleteAllUsers(r.Context())
	if err != nil {
		log.Printf("Error deleting users: %s", err)
		sendError(w, 400, "Failed to delete users")
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}