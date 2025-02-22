package main

import (
	"fmt"
	"net/http"
)

// returns number of file server hits
func (cfg *apiConfig) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	fileServerHitsString := fmt.Sprintf("Hits: %d", cfg.fileServerHits.Load())
		w.Write([]byte(fileServerHitsString))
}