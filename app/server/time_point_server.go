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

	for i := 0; i <= dayRange; i++ {
		today := time.Now().AddDate(0, 0, -i)
		todayY, todayM, todayD := today.Date()
		startOfDay := time.Date(todayY, todayM, todayD, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 0))
		endOfDay := time.Date(todayY, todayM, todayD, 23, 59, 59, 0, time.FixedZone("Asia/Tokyo", 0))

		timePoints, err := models.TimePoints(models.TimePointWhere.PushedAt.GTE(startOfDay), models.TimePointWhere.PushedAt.LTE(endOfDay)).All(r.Context(), t.db)
		if err != nil {
			log.Printf("failed to fetch time points. err: %v", err)
			http.Error(w, fmt.Sprintf("failed to fetch time points. err: %v", err), http.StatusInternalServerError)

			return
		}

		// timePoints から合計時間を集計する
		peer := [][]time.Time{}

		// start, end でペアを組んで 二次元配列に格納する
		for i := 0; i < len(timePoints); i++ {
			tp1 := timePoints[i]

			if i == 0 && tp1.Status == "end" {
				// 最初の要素で, status == "end" だったら00:00 とペアにする
				peer = append(peer, []time.Time{startOfDay, tp1.PushedAt})
				continue
			} else if i == len(timePoints)-1 && tp1.Status == "start" {
				// 最後の要素で, status == "start" だったら23:59 とペアにする
				peer = append(peer, []time.Time{tp1.PushedAt, endOfDay})
				continue
			}

			// 上記以外のパターンは連番でペアを組み、iをさらに1進める(合計 2 進む)
			tp2 := timePoints[i+1]
			peer = append(peer, []time.Time{tp1.PushedAt, tp2.PushedAt})
			i++
		}

		// 作成した peer 稼働時間を算出する
		sum := time.Duration(0)
		for _, p := range peer {
			diff := p[1].Sub(p[0])
			sum += diff
		}
		h := sum / time.Hour
		sum -= h * time.Hour
		m := sum / time.Minute
		resp[domain.Day(fmt.Sprintf("%02d/%02d\n", todayM, todayD))] = domain.TimePointsSum(fmt.Sprintf("%02d:%02d\n", h, m))
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
