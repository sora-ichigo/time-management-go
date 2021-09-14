package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"starter-restapi-golang/app/models"
	"starter-restapi-golang/app/validator"

	_ "github.com/go-sql-driver/mysql"
)

type ContentHandler interface {
	GetContentsHandler(w http.ResponseWriter, r *http.Request)
}

type contentHandlerImpl struct {
	ctx       context.Context
	db        *sql.DB
	validator validator.ContentValidator
}

func NewContentHandler(ctx context.Context, db *sql.DB, v validator.ContentValidator) ContentHandler {
	return &contentHandlerImpl{ctx: ctx, db: db, validator: v}
}

func (u *contentHandlerImpl) GetContentsHandler(w http.ResponseWriter, r *http.Request) {
	contents, err := models.Contents().All(u.ctx, u.db)
	if err != nil {
		log.Fatalf("failed models.Contents(): %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(contents)
	if err != nil {
		log.Printf("failed to encode json. err: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		log.Printf("failed to w.Write(). err: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
