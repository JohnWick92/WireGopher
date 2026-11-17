package handlers

import (
	"net/http"
)

func GetHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, "It's okay for now")
}
