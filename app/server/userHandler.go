package server

import (
	"net/http"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Content-Length", "1");
	w.WriteHeader(http.StatusOK)
}
