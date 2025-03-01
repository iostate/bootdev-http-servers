package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	auth "github.com/iostate/bootdev-http-servers/internal"
	"github.com/iostate/bootdev-http-servers/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	db             *database.Queries
	platform       string
	jwtSecret      string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) storeUserInContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jwtToken, err := auth.GetBearerToken(r.Header)
		if err != nil {
			log.Printf("token is wrong? %s", err)
			respondWithError(w, http.StatusUnauthorized, "Error getting bearer token", err)
			return
		}

		// break down JWT and get the user id
		userID, err := auth.ValidateJWT(jwtToken, cfg.jwtSecret)
		if err != nil {
			log.Printf("token could not be validated")
			respondWithError(w, http.StatusUnauthorized, "Could not validate JWT", err)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		r = r.WithContext(ctx)

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
		log.Printf("Error opening databsae: %v", err)
	}
	dbQueries := database.New(dbConn)

	// Get JWT Secret
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT Secret not set")
	}

	// Add database.Queries to apiCfg
	apiCfg := apiConfig{
		fileServerHits: atomic.Int32{},
		db:             dbQueries,
		platform:       os.Getenv("PLATFORM"),
		jwtSecret:      jwtSecret,
	}

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("."))
	fsHandler := http.StripPrefix("/app", fs)

	// wow, route parameter extraction not allowed with http pkg
	mux.Handle("DELETE /api/chirps/{chirpID}", apiCfg.storeUserInContextMiddleware(http.HandlerFunc(apiCfg.handlerChirpsDelete)))
	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	mux.HandleFunc("GET /api/chirps", apiCfg.handlerGetChirps)
	mux.Handle("PUT /api/users", apiCfg.storeUserInContextMiddleware(http.HandlerFunc(apiCfg.handlerUsersUpdate)))
	mux.HandleFunc("POST /api/refresh", apiCfg.handleRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handleRevoke)
	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirpsById)
	mux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fsHandler))

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	server.ListenAndServe()
}
