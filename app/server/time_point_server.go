package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"time_management_slackapp/app/domain"
	"time_management_slackapp/app/models"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

type TimePointServer interface {
	CreateTimePoint(w http.ResponseWriter, r *http.Request)
	GetTimePointsSumOfDays(w http.ResponseWriter, r *http.Request)
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

	timePoint := models.TimePoint{Status: input.Status}

	// validation
	ok, message := domain.CreateTimePointValidation(timePoint, r.Context(), t.db)
	if !ok {
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	// save DB
	if err := timePoint.Insert(r.Context(), t.db, boil.Infer()); err != nil {
		log.Printf("failed to create time point. err: %v", err)
		http.Error(w, fmt.Sprintf("failed to create time point. err: %v", err), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)

	return
}

func (t timePointServerImpl) GetTimePointsSumOfDays(w http.ResponseWriter, r *http.Request) {
	// 時間取得の日数 0 ~ 6 日前まで
	fetchDays := 0
	fetchDaysQuery, ok := r.URL.Query()["fetch_days"]
	if ok && fetchDaysQuery[0] != "" {
		var err error
		fetchDays, err = strconv.Atoi(fetchDaysQuery[0])
		if err != nil {
			log.Printf("failed to parse day_range err: %v", err)
			http.Error(w, fmt.Sprintf("failed to parse day_range err: %v", err), http.StatusBadRequest)

			return
		}

		if fetchDays > 6 {
			fetchDays = 6
		}
	}
	resp := domain.GetTimePointsSumResponse{}

	// それぞれの日で 稼働時間を算出する
	for i := 0; i <= fetchDays; i++ {
		date := time.Now().AddDate(0, 0, -i)
		_, dateM, dateD := date.Date()

		var err error
		resp[domain.Day(fmt.Sprintf("%02d/%02d", dateM, dateD))], err = domain.CalcTimePointSumOfDay(r.Context(), date, t.db)
		if err != nil {
			log.Printf("failed to fetch time points. err: %v", err)
			http.Error(w, fmt.Sprintf("failed to fetch time points. err: %v", err), http.StatusInternalServerError)

			return
		}
	}

	b, err := json.Marshal(resp)
	if err != nil {
		log.Printf("failed json Marshal err: %v", err)

		return
	}

	if _, err := w.Write(b); err != nil {
		log.Printf("failed to w.Write(). err: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	return
}
