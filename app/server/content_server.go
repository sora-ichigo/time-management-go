package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time_management_slackapp/app/models"
	"time_management_slackapp/app/validator"
	"strconv"

	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
	"github.com/volatiletech/sqlboiler/v4/boil"
)

type ContentServer interface {
	GetContents(w http.ResponseWriter, r *http.Request)
	CreateContent(w http.ResponseWriter, r *http.Request)
	UpdateContent(w http.ResponseWriter, r *http.Request)
	DeleteContent(w http.ResponseWriter, r *http.Request)
}

type contentServerImpl struct {
	ctx       context.Context
	db        *sql.DB
	validator validator.ContentValidator
}

type createContentBody struct {
	Name string
	Text string
}

type updateContentBody createContentBody

func NewContentServer(ctx context.Context, db *sql.DB, v validator.ContentValidator) ContentServer {
	return &contentServerImpl{ctx: ctx, db: db, validator: v}
}

func (u *contentServerImpl) GetContents(w http.ResponseWriter, r *http.Request) {
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

func (u *contentServerImpl) CreateContent(w http.ResponseWriter, r *http.Request) {
	body := createContentBody{}
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

func (u *contentServerImpl) UpdateContent(w http.ResponseWriter, r *http.Request) {
	contentID, err := strconv.ParseUint(chi.URLParam(r, "contentID"), 10, 64)
	if err != nil {
		log.Printf("Invalid URLParam. err: %v", err)
		http.NotFound(w, r)
		return
	}

	body := updateContentBody{}
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

func (u *contentServerImpl) DeleteContent(w http.ResponseWriter, r *http.Request) {
	contentID, err := strconv.ParseUint(chi.URLParam(r, "contentID"), 10, 64)
	if err != nil {
		log.Printf("Invalid URLParam. err: %v", err)
		http.NotFound(w, r)
		return
	}

	content, err := models.FindContent(u.ctx, u.db, uint(contentID))
	if err != nil {
		log.Printf("failed to find content. err: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_, err = content.Delete(u.ctx, u.db)
	if err != nil {
		log.Printf("failed to delete content. err: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
