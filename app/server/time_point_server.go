package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time_management_slackapp/app/domain"
	"time_management_slackapp/app/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

type TimePointServer interface {
	CreateTimePoint(w http.ResponseWriter, r *http.Request)
}

type timePointServerImpl struct {
	db *sql.DB
}

func NewTimePointServer(db *sql.DB) TimePointServer {
	return &timePointServerImpl{db: db}
}

func (t timePointServerImpl) CreateTimePoint(w http.ResponseWriter, r *http.Request) {
	input := domain.CreateTimePointInput{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("Invalid request body. err: %v", err)
		http.Error(w, fmt.Sprintf("Invalid request body. err: %v", err), http.StatusBadRequest)

		return
	}

	// validation
	if input.Status == "" {
		log.Printf("Invalid request body.")
		http.Error(w, "Invalid request body.", http.StatusBadRequest)

		return
	}

	timePoint := models.TimePoint{Status: input.Status}
	fmt.Printf("timePoint: %v", timePoint)

	if err := timePoint.Insert(r.Context(), t.db, boil.Infer()); err != nil {
		log.Printf("failed to create time point. err: %v", err)
		http.Error(w, fmt.Sprintf("failed to create time point. err: %v", err), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)

	return
}
