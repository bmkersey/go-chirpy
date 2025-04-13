package main

import (
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}


func main(){
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
	}
	serverMux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("."))
	serverMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	
	serverMux.HandleFunc("GET /api/healthz", healthHandler)

	serverMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	serverMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	server := &http.Server{
		Addr: ":8080",
		Handler: serverMux,
	}



	

	server.ListenAndServe()



}


