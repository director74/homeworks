package internalhttp

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/director74/homeworks/hw12_13_14_15_calendar/internal/app"
	"github.com/director74/homeworks/hw12_13_14_15_calendar/internal/cfg"
	"github.com/director74/homeworks/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	srv     *http.Server
	host    string
	port    string
	logg    app.Logger
	storage app.Storage
}

type Application interface {
	GetHTTPServerConf() cfg.HTTPServerConf
	GetStorage() app.Storage
}

type Response struct {
	Data  interface{} `json:"data"`
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

type AddResponse struct {
	ID int64 `json:"id"`
}

type ListResponse struct {
	Events []EventRequest `json:"events"`
}

type EventRequest struct {
	ID                   int64   `json:"id,omitempty"`
	Title                *string `json:"title,omitempty"`
	DateStart            *string `json:"dateStart,omitempty"`
	DateEnd              *string `json:"dateEnd,omitempty"`
	Description          *string `json:"description,omitempty"`
	UserID               *int64  `json:"userId,omitempty"`
	NotificationInterval *int32  `json:"notificationInterval,omitempty"`
}

type DeleteRequest struct {
	ID int64 `json:"id,omitempty"`
}

func NewServer(logger app.Logger, app Application) *Server {
	conf := app.GetHTTPServerConf()
	return &Server{
		port:    conf.Port,
		host:    conf.Host,
		logg:    logger,
		storage: app.GetStorage(),
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/add/", s.add)
	mux.HandleFunc("/edit/", s.edit)
	mux.HandleFunc("/delete/", s.delete)
	mux.HandleFunc("/listeventsday/", s.listEventsDay)
	mux.HandleFunc("/listeventsweek/", s.listEventsWeek)
	mux.HandleFunc("/listeventsmonth/", s.listEventsMonth)

	s.srv = &http.Server{
		Addr:         s.host + ":" + s.port,
		Handler:      loggingMiddleware(mux, s.logg),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	s.logg.Info(fmt.Sprintf("starting http server on %s", s.srv.Addr))
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func (s *Server) add(w http.ResponseWriter, r *http.Request) {
	resp := &Response{}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if r.Method != http.MethodPost {
		resp.Error.Message = fmt.Sprintf("method %s not not supported on uri %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusMethodNotAllowed)
		s.writeResponse(w, resp)
		return
	}

	rqEvent := &EventRequest{}
	err := json.NewDecoder(r.Body).Decode(rqEvent)
	if err != nil {
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		s.writeResponse(w, resp)
		return
	}

	id, err := s.storage.Add(s.convertRequestEvent(storage.Event{}, rqEvent))
	if err != nil {
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		s.writeResponse(w, resp)
		return
	}

	resp.Data = AddResponse{ID: id}

	w.WriteHeader(http.StatusOK)
	s.writeResponse(w, resp)
}

func (s *Server) edit(w http.ResponseWriter, r *http.Request) {
	resp := &Response{}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if r.Method != http.MethodPost {
		resp.Error.Message = fmt.Sprintf("method %s not not supported on uri %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusMethodNotAllowed)
		s.writeResponse(w, resp)
		return
	}

	rqEvent := &EventRequest{}
	err := json.NewDecoder(r.Body).Decode(rqEvent)
	if err != nil {
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		s.writeResponse(w, resp)
		return
	}

	foundedEvent, err := s.storage.GetByID(rqEvent.ID)
	if err != nil {
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		s.writeResponse(w, resp)
		return
	}

	err = s.storage.Edit(foundedEvent.ID, s.convertRequestEvent(foundedEvent, rqEvent))
	if err != nil {
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		s.writeResponse(w, resp)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) delete(w http.ResponseWriter, r *http.Request) {
	resp := &Response{}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if r.Method != http.MethodPost {
		resp.Error.Message = fmt.Sprintf("method %s not not supported on uri %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusMethodNotAllowed)
		s.writeResponse(w, resp)
		return
	}

	deleteRq := DeleteRequest{}
	err := json.NewDecoder(r.Body).Decode(&deleteRq)
	if err != nil {
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		s.writeResponse(w, resp)
		return
	}

	err = s.storage.Delete(deleteRq.ID)
	if err != nil {
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		s.writeResponse(w, resp)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) listEventsDay(w http.ResponseWriter, r *http.Request) {
	resp := &Response{}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if r.Method != http.MethodGet {
		resp.Error.Message = fmt.Sprintf("method %s not not supported on uri %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusMethodNotAllowed)
		s.writeResponse(w, resp)
		return
	}

	dayRq := r.URL.Query().Get("date")
	if dayRq == "" {
		resp.Error.Message = "date not specified"
		w.WriteHeader(http.StatusBadRequest)
		s.writeResponse(w, resp)
		return
	}

	events, err := s.storage.ListEventsDay(dayRq)
	if err != nil {
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		s.writeResponse(w, resp)
		return
	}

	responseEvents := make([]EventRequest, len(events))
	for i, event := range events {
		responseEvents[i] = s.convertResponseEvent(event)
	}
	resp.Data = ListResponse{Events: responseEvents}

	w.WriteHeader(http.StatusOK)
	s.writeResponse(w, resp)
}

func (s *Server) listEventsWeek(w http.ResponseWriter, r *http.Request) {
	resp := &Response{}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if r.Method != http.MethodGet {
		resp.Error.Message = fmt.Sprintf("method %s not not supported on uri %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusMethodNotAllowed)
		s.writeResponse(w, resp)
		return
	}

	dayRq := r.URL.Query().Get("weekBeginDate")
	if dayRq == "" {
		resp.Error.Message = "begin week date not specified"
		w.WriteHeader(http.StatusBadRequest)
		s.writeResponse(w, resp)
		return
	}

	events, err := s.storage.ListEventsWeek(dayRq)
	if err != nil {
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		s.writeResponse(w, resp)
		return
	}

	responseEvents := make([]EventRequest, len(events))
	for i, event := range events {
		responseEvents[i] = s.convertResponseEvent(event)
	}
	resp.Data = ListResponse{Events: responseEvents}

	w.WriteHeader(http.StatusOK)
	s.writeResponse(w, resp)
}

func (s *Server) listEventsMonth(w http.ResponseWriter, r *http.Request) {
	resp := &Response{}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if r.Method != http.MethodGet {
		resp.Error.Message = fmt.Sprintf("method %s not not supported on uri %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusMethodNotAllowed)
		s.writeResponse(w, resp)
		return
	}

	dayRq := r.URL.Query().Get("monthBeginDate")
	if dayRq == "" {
		resp.Error.Message = "begin month date not specified"
		w.WriteHeader(http.StatusBadRequest)
		s.writeResponse(w, resp)
		return
	}

	events, err := s.storage.ListEventsMonth(dayRq)
	if err != nil {
		resp.Error.Message = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		s.writeResponse(w, resp)
		return
	}

	responseEvents := make([]EventRequest, len(events))
	for i, event := range events {
		responseEvents[i] = s.convertResponseEvent(event)
	}
	resp.Data = ListResponse{Events: responseEvents}

	w.WriteHeader(http.StatusOK)
	s.writeResponse(w, resp)
}

func (s *Server) writeResponse(w http.ResponseWriter, resp *Response) {
	resBuf, err := json.Marshal(resp)
	if err != nil {
		log.Printf("response marshal error: %s", err)
	}
	_, err = w.Write(resBuf)
	if err != nil {
		log.Printf("response marshal error: %s", err)
	}
}

func (s *Server) convertRequestEvent(storageEvent storage.Event, event *EventRequest) storage.Event {
	if event.Title != nil {
		storageEvent.Title = *event.Title
	}

	if event.DateStart != nil {
		newDateStart, _ := time.ParseInLocation(storage.DateTimeFormatISO, *event.DateStart, time.Local)
		storageEvent.DateStart = newDateStart
	}

	if event.DateEnd != nil {
		newDateEnd, _ := time.ParseInLocation(storage.DateTimeFormatISO, *event.DateEnd, time.Local)
		storageEvent.DateEnd = newDateEnd
	}

	if event.Description != nil {
		storageEvent.Description = sql.NullString{String: *event.Description, Valid: true}
	}

	if event.UserID != nil {
		storageEvent.UserID = *event.UserID
	}

	if event.NotificationInterval != nil {
		storageEvent.NotificationInterval = sql.NullInt32{Int32: *event.NotificationInterval, Valid: true}
	}

	return storageEvent
}

func (s *Server) convertResponseEvent(storageEvent storage.Event) EventRequest {
	dateStart := storageEvent.DateStart.Format(storage.DateTimeFormatISO)
	dateEnd := storageEvent.DateEnd.Format(storage.DateTimeFormatISO)
	return EventRequest{
		ID:                   storageEvent.ID,
		Title:                &storageEvent.Title,
		DateStart:            &dateStart,
		DateEnd:              &dateEnd,
		Description:          &storageEvent.Description.String,
		UserID:               &storageEvent.UserID,
		NotificationInterval: &storageEvent.NotificationInterval.Int32,
	}
}
