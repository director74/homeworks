package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/director74/homeworks/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	events map[int64]storage.Event
	mu     sync.RWMutex
	lastID int64
}

func New() *Storage {
	return &Storage{
		events: make(map[int64]storage.Event),
	}
}

func (s *Storage) Add(ctx context.Context, event storage.Event) (int64, error) {
	if err := s.validate(event); err != nil {
		return 0, err
	}
	s.mu.Lock()
	newIndex := s.lastID + 1
	event.ID = newIndex
	s.events[newIndex] = event
	s.lastID = newIndex
	s.mu.Unlock()

	return newIndex, nil
}

func (s *Storage) Edit(ctx context.Context, i int64, event storage.Event) error {
	if err := s.validate(event); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[i]; !ok {
		return fmt.Errorf("element %v doesnt exist", i)
	}
	s.events[i] = event

	return nil
}

func (s *Storage) Delete(ctx context.Context, i int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[i]; !ok {
		return fmt.Errorf("element %v doesnt exist", i)
	}
	delete(s.events, i)

	return nil
}

func (s *Storage) GetByID(id int64) (storage.Event, error) {
	s.mu.RLock()
	event, ok := s.events[id]
	s.mu.RUnlock()

	if ok {
		return event, nil
	}

	return storage.Event{}, storage.ErrEventNotFound
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

	return s.listByDates(dt, dt2)
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

	return s.listByDates(dt, dt2)
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

	return s.listByDates(dt, dt2)
}

func (s *Storage) listByDates(dt time.Time, dt2 time.Time) ([]storage.Event, error) {
	events := make([]storage.Event, 0)

	s.mu.RLock()
	for _, event := range s.events {
		if (event.DateStart.After(dt) && event.DateStart.Before(dt2)) ||
			(event.DateEnd.After(dt) && event.DateEnd.Before(dt2)) ||
			event.DateStart.Equal(dt) ||
			event.DateStart.Equal(dt2) ||
			event.DateEnd.Equal(dt) ||
			event.DateEnd.Equal(dt2) {
			events = append(events, event)
		}
	}
	s.mu.RUnlock()

	return events, nil
}

func (s *Storage) validate(event storage.Event) error {
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

func (s *Storage) count() int {
	return len(s.events)
}
