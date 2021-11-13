package domain

import (
	"context"
	"database/sql"
	"errors"
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
type TimePointsSum string

// 1日あたりの稼働時間
type GetTimePointsSumResponse map[Day]TimePointsSum

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
		return false, "failed find timePoint record"
	}

	if prevTimePoint == nil && t.Status == "end" || (prevTimePoint != nil && prevTimePoint.Status == t.Status) {
		return false, "invalid status peer"
	}
	return true, ""
}
