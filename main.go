package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/iostate/bootdev-http-servers/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db 	*database.Queries
	platform string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	
	// Start SQL
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Errorf("Error opening databsae: %w", err)
	}
	dbQueries := database.New(dbConn)

	// Add database.Queries to apiCfg
	apiCfg := apiConfig{
		fileServerHits: atomic.Int32{},
		db: dbQueries,
		platform: os.Getenv("PLATFORM"),
	}
	

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("."))
	fsHandler :=  http.StripPrefix("/app", fs)
	
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpsById)
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fsHandler))

	server := &http.Server{
		Handler: mux,
		Addr: ":8080",
	}

	server.ListenAndServe()
}