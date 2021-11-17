package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"time_management_slackapp/app/models"

	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

type CreateTimePointInput struct {
	Status string
}

// 何日前から取得するか
type GetTimePointsSumInput struct {
	DaysAgo int
}

// e.g. 01/01, 12/10
type Day string

// 何時間稼働したか
// e.g. "03:20"
type TimePointSum string

// 1日あたりの稼働時間
type GetTimePointsSumResponse map[Day]TimePointSum

func CreateTimePointValidation(t models.TimePoint, ctx context.Context, db *sql.DB) (ok bool, message string) {
	if t.Status == "" {
		return false, "status is required"
	}

	if !(t.Status == "start" || t.Status == "end") {
		return false, "status must be `end` or `start`"
	}

	// start と end の対応があっているかの確認
	// 最新の一件を取得
	prevTimePoint, err := models.TimePoints(qm.OrderBy("id DESC")).One(ctx, db)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, "failed to find timePoint record"
	}

	if prevTimePoint == nil && t.Status == "end" || (prevTimePoint != nil && prevTimePoint.Status == t.Status) {
		return false, "invalid status peer"
	}
	return true, ""
}

// 特定の日の稼働時間を算出する
func CalcTimePointSumOfDay(ctx context.Context, date time.Time, db *sql.DB) (TimePointSum, error) {
	year, month, day := date.Date()
	startOfDay := time.Date(year, month, day, 0, 0, 0, 0, time.FixedZone("Asia/Tokyo", 0))
	endOfDay := time.Date(year, month, day, 23, 59, 59, 0, time.FixedZone("Asia/Tokyo", 0))

	timePoints, err := models.TimePoints(models.TimePointWhere.PushedAt.GTE(startOfDay), models.TimePointWhere.PushedAt.LTE(endOfDay)).All(ctx, db)
	if err != nil {
		return TimePointSum(""), err
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
	return TimePointSum(fmt.Sprintf("%02d:%02d", h, m)), nil
}
