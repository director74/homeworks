package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"sync"
	"syscall"

	"github.com/director74/homeworks/hw12_13_14_15_calendar/internal/app"
	"github.com/director74/homeworks/hw12_13_14_15_calendar/internal/cfg"
	"github.com/director74/homeworks/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/director74/homeworks/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/director74/homeworks/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/director74/homeworks/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/director74/homeworks/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./config.yml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := cfg.NewConfig()
	err := config.Parse(configFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	ctx := context.Background()
	logg := logger.New(config.GetLoggerConf().Level)

	var storage app.Storage

	switch config.GetAppConf().StorageType {
	case cfg.SQLStorage:
		dbStorage, err := sqlstorage.New(ctx, config.GetDBConf())
		defer dbStorage.Close()
		storage = dbStorage
		if err != nil {
			logg.Error(err.Error())
			return
		}
	case cfg.MemoryStorage:
		fallthrough
	default:
		storage = memorystorage.New()
	}
	app.New(logg, storage, config)

	httpServer := internalhttp.NewServer(logg, storage, config.GetServersConf().HTTP)
	grpcServer := internalgrpc.NewServer(logg, storage, config.GetServersConf().GRPC)

	ctx, cancel := signal.NotifyContext(ctx,
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	wg := &sync.WaitGroup{}

	go func() {
		<-ctx.Done()
		if err := httpServer.Stop(ctx); err != nil {
			logg.Error("http server stopped: " + err.Error())
		}
	}()

	go func() {
		<-ctx.Done()
		grpcServer.Stop()
	}()

	logg.Info("calendar is running...")

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.Start(); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := grpcServer.Start(); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
		}
	}()

	wg.Wait()
}
