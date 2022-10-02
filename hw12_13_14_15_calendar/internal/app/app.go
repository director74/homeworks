package app

import (
	"context"

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
	Add(context.Context, storage.Event) (int64, error)
	Edit(context.Context, int64, storage.Event) error
	Delete(context.Context, int64) error
	ListEventsDay(ctx context.Context, date string) ([]storage.Event, error)
	ListEventsWeek(ctx context.Context, weekBeginDate string) ([]storage.Event, error)
	ListEventsMonth(ctx context.Context, monthBeginDate string) ([]storage.Event, error)
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
