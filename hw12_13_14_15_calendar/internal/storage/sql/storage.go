package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/cfg"
	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/jackc/pgx/stdlib" // linter said that need comment
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

// TODO Как здесь использовать контекст?

func (s *Storage) Close(ctx context.Context) error {
	if err := s.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	return nil
}

func (s *Storage) Add(event storage.Event) (int64, error) {
	var returnedID int64

	err := s.validate(event)
	if err != nil {
		return 0, err
	}

	res, err := s.conn.NamedQuery(`
		INSERT INTO events ("Title", "DateStart", "DateEnd", "Description", "UserID", "NotificationInterval") 
		VALUES (:Title, :DateStart, :DateEnd, :Description, :UserID, :NotificationInterval)
		RETURNING "ID"`,
		event,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to add record: %w", err)
	}
	defer res.Close()

	for res.Next() {
		err = res.Scan(&returnedID)
		if err != nil {
			return 0, fmt.Errorf("failed to scan record: %w", err)
		}
	}

	return returnedID, nil
}

func (s *Storage) Edit(i int64, event storage.Event) error {
	err := s.validate(event)
	if err != nil {
		return err
	}

	event.ID = i
	_, err = s.conn.NamedExec(`
		UPDATE events 
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

func (s *Storage) Delete(i int64) error {
	_, err := s.conn.Exec(`DELETE FROM events WHERE "ID" = $1`, i)
	if err != nil {
		return fmt.Errorf("failed to delete record: %w", err)
	}

	return nil
}

func (s *Storage) ListEventsDay(date string) ([]storage.Event, error) {
	dt, err := time.ParseInLocation(storage.DateFormatISO, date, time.Local)
	if err != nil {
		return nil, fmt.Errorf("wrong date format: %w", err)
	}

	dt2 := dt.Add(time.Hour * 23).Add(time.Minute * 59).Add(time.Second * 59)

	return s.listByDates(dt, dt2)
}

func (s *Storage) ListEventsWeek(weekBeginDate string) ([]storage.Event, error) {
	dt, err := time.ParseInLocation(storage.DateFormatISO, weekBeginDate, time.Local)
	if err != nil {
		return nil, fmt.Errorf("wrong date format: %w", err)
	}

	dt2 := dt.AddDate(0, 0, 6).Add(time.Hour * 23).Add(time.Minute * 59).Add(time.Second * 59)

	return s.listByDates(dt, dt2)
}

func (s *Storage) ListEventsMonth(monthBeginDate string) ([]storage.Event, error) {
	dt, err := time.ParseInLocation(storage.DateFormatISO, monthBeginDate, time.Local)
	if err != nil {
		return nil, fmt.Errorf("wrong date format: %w", err)
	}

	dt2 := dt.AddDate(0, 1, -dt.Day())

	return s.listByDates(dt, dt2)
}

func (s *Storage) listByDates(dt time.Time, dt2 time.Time) ([]storage.Event, error) {
	events := make([]storage.Event, 0)

	err := s.conn.Select(
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

func (s *Storage) validate(event storage.Event) error {
	if event.Title == "" {
		return storage.ErrWrongTitle
	}
	if event.DateStart.Before(time.Now()) {
		return storage.ErrWrongDateStart
	}
	if event.DateEnd.Before(event.DateStart) {
		return storage.ErrWrongDateEnd
	}
	if s.checkBusy(event.DateStart, event.DateEnd, event.ID) {
		return storage.ErrDateBusy
	}
	return nil
}

func (s *Storage) checkBusy(dateStart time.Time, dateEnd time.Time, id int64) bool {
	events, _ := s.listByDates(dateStart, dateEnd)
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
