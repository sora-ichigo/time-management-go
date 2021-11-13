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
	GetTimePointsSum(w http.ResponseWriter, r *http.Request)
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

func (t timePointServerImpl) GetTimePointsSum(w http.ResponseWriter, r *http.Request) {
	// 時間取得の日数 0 ~ 6 日前まで
	dayRange := 0
	dayRangeQuery, ok := r.URL.Query()["day_range"]
	if ok {
		var err error
		dayRange, err = strconv.Atoi(dayRangeQuery[0])
		if err != nil {
			log.Printf("failed to parse day_range err: %v", err)
			http.Error(w, fmt.Sprintf("failed to parse day_range err: %v", err), http.StatusBadRequest)

			return
		}

		if dayRange > 6 {
			dayRange = 6
		}
	}
	resp := domain.GetTimePointsSumResponse{}

	// それぞれの日で 稼働時間を算出する
	for i := 0; i <= dayRange; i++ {
		date := time.Now().AddDate(0, 0, -i)
		_, dateM, dateD := date.Date()

		var err error
		resp[domain.Day(fmt.Sprintf("%02d/%02d\n", dateM, dateD))], err = domain.CalcTimePointSum(date, r.Context(), t.db)
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
