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
	ErrWrongUserID    = errors.New("event user can't be empty")
	ErrDeleteNotFound = errors.New("event for deletion not found")
	ErrEventNotFound  = errors.New("event not found")
	ErrEmptyDate      = errors.New("date is empty")
	ErrSameDates      = errors.New("begin and end time are the same")
)

const (
	DateFormatISO     = "2006-01-02"
	DateTimeFormatISO = "2006-01-02 15:04:05"
	LayoutLog         = "01/Jan/2006 03:04:05 -0700"
)

type Event struct {
	ID                   int64          `db:"ID" fake:"{number:1,9223372036854775807}" json:"id"`
	Title                string         `db:"Title" fake:"{sentence:10}" json:"title"`
	DateStart            time.Time      `db:"DateStart" json:"dateStart"`
	DateEnd              time.Time      `db:"DateEnd" json:"dateEnd"`
	Description          sql.NullString `db:"Description" json:"description"`
	UserID               int64          `db:"UserID" fake:"{number:1,9223372036854775807}" json:"userId"`
	NotificationInterval sql.NullInt32  `db:"NotificationInterval" json:"notificationInterval"`
}
