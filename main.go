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
	dbQueries *database.Queries
}


func main(){
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	dbQueries := database.New(db)
	if err != nil {
		fmt.Println("Error connecting to DB")
		return
	}

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries: dbQueries, 
	}

	serverMux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("."))
	serverMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))

	serverMux.HandleFunc("GET /api/healthz", healthHandler)
	serverMux.HandleFunc("POST /api/validate_chirp", apiCfg.handlerValidateChirp)

	serverMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	serverMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	server := &http.Server{
		Addr: ":8080",
		Handler: serverMux,
	}

	server.ListenAndServe()

}


