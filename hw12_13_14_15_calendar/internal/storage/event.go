package storage

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrDateBusy       = errors.New("selected time already reserved by another event")
	ErrWrongDateStart = errors.New("event begin date can't be in past")
	ErrWrongDateEnd   = errors.New("event finish can't be same or earlier then start")
	ErrWrongTitle     = errors.New("event title can't be empty")
)

const (
	DateFormatISO = "2006-01-02"
	LayoutLog     = "01/Jan/2006 03:04:05 -0700"
)

type Event struct {
	ID                   int64          `db:"ID" fake:"{number:1,9223372036854775807}"`
	Title                string         `db:"Title" fake:"{sentence:10}"`
	DateStart            time.Time      `db:"DateStart"`
	DateEnd              time.Time      `db:"DateEnd"`
	Description          sql.NullString `db:"Description"`
	UserID               int64          `db:"UserID" fake:"{number:1,9223372036854775807}"`
	NotificationInterval sql.NullInt16  `db:"NotificationInterval"`
}
