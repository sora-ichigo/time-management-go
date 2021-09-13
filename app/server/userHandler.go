package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"starter-restapi-golang/app/models"

	_ "github.com/go-sql-driver/mysql"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := models.Users().AllG(context.Context)
	if err != nil {
		log.Fatalf("failed models.Users(): %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(users)
	if err != nil {
		log.Printf("failed to encode json. err: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
