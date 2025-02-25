package main

import (
	"fmt"
	"net/http"
)

// returns HTML with number of file server hits
func (cfg *apiConfig) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type", "text/html")
	htmlToReturn := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileServerHits.Load())

		w.Write([]byte(htmlToReturn))
}