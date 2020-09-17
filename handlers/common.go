package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Deps contains dependencies of the http handlers
type Deps struct {
	DB *pgxpool.Pool
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}
