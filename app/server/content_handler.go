package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"starter-restapi-golang/app/models"
	"starter-restapi-golang/app/validator"
	"strconv"

	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type ContentHandler interface {
	GetContentsHandler(w http.ResponseWriter, r *http.Request)
	PostContentsHandler(w http.ResponseWriter, r *http.Request)
	PutContentHandler(w http.ResponseWriter, r *http.Request)
}

type contentHandlerImpl struct {
	ctx       context.Context
	db        *sql.DB
	validator validator.ContentValidator
}

type postContentBody struct {
	Name string
	Text string
}

type putContentBody postContentBody

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

func (u *contentHandlerImpl) PostContentsHandler(w http.ResponseWriter, r *http.Request) {
	body := postContentBody{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("Invalid request body. err: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	content := models.Content{Name: body.Name, Text: body.Text}
	err := u.validator.Set(content).Valid()
	if err != nil {
		log.Printf("Invalid request body. err: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = content.Insert(u.ctx, u.db, boil.Infer())
	if err != nil {
		log.Printf("failed to create content. err: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (u *contentHandlerImpl) PutContentHandler(w http.ResponseWriter, r *http.Request) {
	contentID, err := strconv.ParseUint(chi.URLParam(r, "contentID"), 10, 64)
	if err != nil {
		log.Printf("Invalid URLParam. err: %v", err)
		http.NotFound(w, r)
		return
	}

	body := putContentBody{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Printf("Invalid request body. err: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	content, err := models.FindContent(u.ctx, u.db, uint(contentID))
	if err != nil {
		log.Printf("failed to find content. err: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	content.Name = body.Name
	content.Text = body.Text

	err = u.validator.Set(*content).Valid()
	if err != nil {
		log.Printf("Invalid request body. err: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = content.Update(u.ctx, u.db, boil.Infer())
	if err != nil {
		log.Printf("failed to update content. err: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
