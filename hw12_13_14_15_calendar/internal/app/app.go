package app

import (
	"context"

	"github.com/director74/homeworks/hw12_13_14_15_calendar/internal/cfg"
	"github.com/director74/homeworks/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
	config  cfg.Config
}

type Logger interface {
	Info(string)
	Error(string)
	Warn(string)
	Debug(string)
}

type Storage interface {
	Add(storage.Event) (int64, error)
	Edit(int64, storage.Event) error
	Delete(int64) error
	ListEventsDay(date string) ([]storage.Event, error)
	ListEventsWeek(weekBeginDate string) ([]storage.Event, error)
	ListEventsMonth(monthBeginDate string) ([]storage.Event, error)
}

func New(logger Logger, storage Storage, config cfg.Config) *App {
	return &App{
		logger:  logger,
		storage: storage,
		config:  config,
	}
}

func (a *App) GetServerConf() cfg.ServerConf {
	return a.config.Server
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
