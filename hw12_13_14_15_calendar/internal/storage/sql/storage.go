package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/director74/homeworks/hw12_13_14_15_calendar/internal/cfg"
	"github.com/director74/homeworks/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	conn *sqlx.DB
}

func New(ctx context.Context, params cfg.DatabaseConf) (*Storage, error) {
	storage := &Storage{}
	err := storage.Connect(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("db storage problem: %w", err)
	}
	return storage, nil
}

func (s *Storage) Connect(ctx context.Context, params cfg.DatabaseConf) (err error) {
	dsn := fmt.Sprintf("user=%s dbname=postgres sslmode=disable password=%s", params.User, params.Password)
	s.conn, err = sqlx.ConnectContext(ctx, "pgx", dsn)

	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}
	return nil
}

func (s *Storage) Close() error {
	if err := s.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	return nil
}

func (s *Storage) Add(ctx context.Context, event storage.Event) (int64, error) {
	var returnedID int64

	err := s.validate(ctx, event)
	if err != nil {
		return 0, err
	}

	returnedID = 0

	err = s.conn.GetContext(
		ctx,
		&returnedID,
		`INSERT INTO events ("Title", "DateStart", "DateEnd", "Description", "UserID", "NotificationInterval")
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING "ID"`,
		event.Title,
		event.DateStart,
		event.DateEnd,
		event.Description,
		event.UserID,
		event.NotificationInterval)
	if err != nil {
		return 0, fmt.Errorf("failed to add record: %w", err)
	}

	return returnedID, nil
}

func (s *Storage) Edit(ctx context.Context, i int64, event storage.Event) error {
	err := s.validate(ctx, event)
	if err != nil {
		return err
	}

	event.ID = i
	_, err = s.conn.NamedExecContext(ctx,
		`UPDATE events 
		SET ("Title", "DateStart", "DateEnd", "Description", "UserID", "NotificationInterval") = 
		    (:Title, :DateStart, :DateEnd, :Description, :UserID, :NotificationInterval)
		WHERE "ID" = :ID`,
		event,
	)
	if err != nil {
		return fmt.Errorf("failed to update record: %w", err)
	}

	return nil
}

func (s *Storage) Delete(ctx context.Context, i int64) error {
	deletedCount := 0

	err := s.conn.GetContext(
		ctx,
		&deletedCount,
		`WITH deleted AS 
    	(DELETE FROM events WHERE "ID" = $1 RETURNING *) SELECT count(*) FROM deleted`,
		i,
	)
	if err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}

	if deletedCount < 1 {
		return storage.ErrDeleteNotFound
	}

	return nil
}

func (s *Storage) GetByID(id int64) (storage.Event, error) {
	event := storage.Event{}
	err := s.conn.Get(&event, `SELECT * FROM events WHERE "ID"=$1`, id)
	if err != nil {
		return storage.Event{}, fmt.Errorf("failed to get event: %w", err)
	}
	if event.ID == 0 {
		return storage.Event{}, storage.ErrEventNotFound
	}

	return event, nil
}

func (s *Storage) ListEventsDay(ctx context.Context, date string) ([]storage.Event, error) {
	if date == "" {
		return nil, storage.ErrEmptyDate
	}
	dt, err := time.ParseInLocation(storage.DateFormatISO, date, time.Local)
	if err != nil {
		return nil, fmt.Errorf("wrong date format: %w", err)
	}

	dt2 := dt.Add(time.Hour * 23).Add(time.Minute * 59).Add(time.Second * 59)

	return s.listByDates(ctx, dt, dt2)
}

func (s *Storage) ListEventsWeek(ctx context.Context, weekBeginDate string) ([]storage.Event, error) {
	if weekBeginDate == "" {
		return nil, storage.ErrEmptyDate
	}
	dt, err := time.ParseInLocation(storage.DateFormatISO, weekBeginDate, time.Local)
	if err != nil {
		return nil, fmt.Errorf("wrong date format: %w", err)
	}

	dt2 := dt.AddDate(0, 0, 6).Add(time.Hour * 23).Add(time.Minute * 59).Add(time.Second * 59)

	return s.listByDates(ctx, dt, dt2)
}

func (s *Storage) ListEventsMonth(ctx context.Context, monthBeginDate string) ([]storage.Event, error) {
	if monthBeginDate == "" {
		return nil, storage.ErrEmptyDate
	}
	dt, err := time.ParseInLocation(storage.DateFormatISO, monthBeginDate, time.Local)
	if err != nil {
		return nil, fmt.Errorf("wrong date format: %w", err)
	}

	dt2 := dt.AddDate(0, 1, -dt.Day())

	return s.listByDates(ctx, dt, dt2)
}

func (s *Storage) listByDates(ctx context.Context, dt time.Time, dt2 time.Time) ([]storage.Event, error) {
	events := make([]storage.Event, 0)

	err := s.conn.SelectContext(
		ctx,
		&events,
		`SELECT * FROM events WHERE "DateStart" BETWEEN $1 AND $2 OR "DateEnd" BETWEEN $1 AND $2 ORDER BY "ID" ASC`,
		dt,
		dt2,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to select events day: %w", err)
	}

	return events, nil
}

func (s *Storage) validate(ctx context.Context, event storage.Event) error {
	if event.Title == "" {
		return storage.ErrWrongTitle
	}
	if event.UserID == 0 {
		return storage.ErrWrongUserID
	}
	if event.DateStart.Before(time.Now()) {
		return storage.ErrWrongDateStart
	}
	if event.DateEnd.Before(event.DateStart) {
		return storage.ErrWrongDateEnd
	}
	if event.DateStart == event.DateEnd {
		return storage.ErrSameDates
	}
	if s.checkBusy(ctx, event.DateStart, event.DateEnd, event.ID) {
		return storage.ErrDateBusy
	}
	return nil
}

func (s *Storage) checkBusy(ctx context.Context, dateStart time.Time, dateEnd time.Time, id int64) bool {
	events, _ := s.listByDates(ctx, dateStart, dateEnd)
	if len(events) > 0 {
		for _, founded := range events {
			if founded.ID != id {
				return true
			}
		}
	}

	return false
}

/*func (s *Storage) count() (int, error) {
	result := []int{}
	err := s.conn.Select(&result, "SELECT COUNT(*) cnt FROM events")
	if err != nil {
		return 0, err
	}

	return result[0], nil
}*/
