package app

import (
	"github.com/director74/homeworks/hw12_13_14_15_calendar/internal/cfg"
	"github.com/director74/homeworks/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
	config  cfg.Configurable
}

type Logger interface {
	Info(string)
	Infof(msg string, args ...interface{})
	Error(string)
	Warn(string)
	Debug(string)
}

//go:generate mockgen -source=./internal/app/app.go --destination=./test/mocks/app/app.go
type Storage interface {
	Add(storage.Event) (int64, error)
	Edit(int64, storage.Event) error
	Delete(int64) error
	ListEventsDay(date string) ([]storage.Event, error)
	ListEventsWeek(weekBeginDate string) ([]storage.Event, error)
	ListEventsMonth(monthBeginDate string) ([]storage.Event, error)
	GetByID(int64) (storage.Event, error)
}

func New(logger Logger, storage Storage, config cfg.Configurable) *App {
	return &App{
		logger:  logger,
		storage: storage,
		config:  config,
	}
}

func (a *App) GetHTTPServerConf() cfg.HTTPServerConf {
	return a.config.GetServersConf().HTTP
}

func (a *App) GetGRPCServerConf() cfg.GRPCServerConf {
	return a.config.GetServersConf().GRPC
}

func (a *App) GetStorage() Storage {
	return a.storage
}
