package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/bmkersey/go-chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
	jwtSecret      string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Error connecting to DB")
		return
	}

	dbQueries := database.New(db)
	platform := os.Getenv("PLATFORM")
	jwtSecret := os.Getenv("JWT_SECRET")

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      dbQueries,
		platform:       platform,
		jwtSecret:      jwtSecret,
	}

	serverMux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("."))
	serverMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))

	serverMux.HandleFunc("GET /api/healthz", healthHandler)
	serverMux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	serverMux.HandleFunc("POST /api/refresh", apiCfg.handlerRefreshToken)
	serverMux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeRefreshToken)

	serverMux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	serverMux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	serverMux.HandleFunc("POST /api/validate_chirp", apiCfg.handlerValidateChirp)
	serverMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetSingleChirp)
	serverMux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)

	serverMux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	serverMux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)

	serverMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	serverMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	server := &http.Server{
		Addr:    ":8080",
		Handler: serverMux,
	}

	server.ListenAndServe()

}
