package main

import (
	"net/http"

	"github.com/bmkersey/go-chirpy/internal/auth"
)

func (c *apiConfig) handlerRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		sendError(w, http.StatusBadRequest, "Couldn't find token")
		return
	}

	_, err = c.dbQueries.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Couldn't revoke session")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
