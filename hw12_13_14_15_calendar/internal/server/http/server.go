package internalhttp

import (
	"context"
	"net/http"
	"time"

	"github.com/director74/homeworks/hw12_13_14_15_calendar/internal/cfg"
)

type Server struct {
	srv  *http.Server
	host string
	port string
	logg Logger
}

type Logger interface {
	Info(string)
	Infof(msg string, args ...interface{})
	Error(string)
	Warn(string)
	Debug(string)
}

type Application interface {
	GetServerConf() cfg.ServerConf
}

func NewServer(logger Logger, app Application) *Server {
	conf := app.GetServerConf()
	return &Server{
		port: conf.Port,
		host: conf.Host,
		logg: logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", s.hello)

	s.srv = &http.Server{
		Addr:         s.host + ":" + s.port,
		Handler:      loggingMiddleware(mux, s.logg),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := s.srv.ListenAndServe(); err != nil {
		return err
	}
	// TODO Как здесь использовать контекст?
	<-ctx.Done()
	return nil
}

// TODO Как здесь использовать контекст?

func (s *Server) Stop(ctx context.Context) error {
	s.srv.Close()
	return nil
}

func (s *Server) hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("World!"))
}
