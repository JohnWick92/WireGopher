package handlers

import (
	"encoding/json"
	"net/http"
)

// writeJSON is a shared helper function to reduce
// boilerplate converting to JSON
func writeJSON(w http.ResponseWriter, status int, msg any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(msg)
}
