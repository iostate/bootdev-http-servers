package main

import (
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) resetToZero(w http.ResponseWriter, r *http.Request) {
	// return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		// reset to zero
		cfg.fileServerHits.Store(0)
		w.Write([]byte("Counter reset to 0"))
}

func main() {
	apiCfg := apiConfig{}

	mux := http.NewServeMux()
	
	// file server, strip the fucking prefix
	fs := http.FileServer(http.Dir("."))
	fsHandler :=  http.StripPrefix("/app", fs)
	
	mux.HandleFunc("POST /api/validate_chirp", ValidateChirpHandler)
	mux.HandleFunc("GET /admin/metrics", apiCfg.MetricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetToZero)
	mux.HandleFunc("GET /api/healthz", ReadinessHandler)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fsHandler))

	server := &http.Server{
		Handler: mux,
		Addr: ":8080",
	}

	server.ListenAndServe()
}