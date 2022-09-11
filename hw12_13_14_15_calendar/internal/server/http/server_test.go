package internalhttp

import (
	"bufio"
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	gofake "github.com/brianvoe/gofakeit/v6"
	"github.com/director74/homeworks/hw12_13_14_15_calendar/internal/app"
	"github.com/director74/homeworks/hw12_13_14_15_calendar/internal/cfg"
	"github.com/director74/homeworks/hw12_13_14_15_calendar/internal/storage"
	mock_app "github.com/director74/homeworks/hw12_13_14_15_calendar/test/mocks/app"
	mock_cfg "github.com/director74/homeworks/hw12_13_14_15_calendar/test/mocks/cfg"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestNewServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLog := mock_app.NewMockLogger(ctrl)
	mockStorage := mock_app.NewMockStorage(ctrl)
	mockConfig := mock_cfg.NewMockConfigurable(ctrl)

	mockConfig.EXPECT().GetServersConf().Return(cfg.ServersConf{GRPC: cfg.GRPCServerConf{}, HTTP: cfg.HTTPServerConf{}})

	apl := app.New(mockLog, mockStorage, mockConfig)
	s := NewServer(mockLog, apl)

	t.Run("Check success add", func(t *testing.T) {
		element, errPrepare := prepareEventRequest(10)
		require.NoError(t, errPrepare)
		jsonEvent, errMarshal := json.Marshal(&element)
		require.NoError(t, errMarshal)

		mockStorage.EXPECT().Add(s.convertRequestEvent(storage.Event{}, &element)).Return(int64(1), nil).Times(1)

		req := httptest.NewRequest(http.MethodPost, "/add/", bytes.NewReader(jsonEvent))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("accept", "application/json")
		w := httptest.NewRecorder()
		s.add(w, req)

		result := w.Result()
		defer result.Body.Close()

		response := &Response{}
		response.Data = &AddResponse{}

		errDecode := json.NewDecoder(result.Body).Decode(response)
		require.NoError(t, errDecode)

		addResp, ok := response.Data.(*AddResponse)
		if !ok {
			addResp = &AddResponse{}
		}

		require.Equal(t, "", response.Error.Message)
		require.Equal(t, addResp.ID, int64(1))
		require.Equal(t, http.StatusOK, result.StatusCode)
	})

	t.Run("Check fail add", func(t *testing.T) {
		mockStorage.EXPECT().Add(storage.Event{}).Return(int64(0), storage.ErrWrongTitle).Times(1)

		jsonEvent, errMarshal := json.Marshal(&EventRequest{})
		require.NoError(t, errMarshal)

		req := httptest.NewRequest(http.MethodPost, "/add/", bytes.NewReader(jsonEvent))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("accept", "application/json")
		w := httptest.NewRecorder()
		s.add(w, req)

		result := w.Result()
		defer result.Body.Close()

		response := &Response{}

		errDecode := json.NewDecoder(result.Body).Decode(response)
		require.NoError(t, errDecode)

		require.Contains(t, response.Error.Message, storage.ErrWrongTitle.Error())
		require.Equal(t, response.Data, nil)
		require.Equal(t, http.StatusBadRequest, result.StatusCode)
	})

	t.Run("Check edit", func(t *testing.T) {
		sEvent, err := prepareStorageEvent("2022-10-01", 1)
		require.NoError(t, err)
		updatedEvent := sEvent
		mockStorage.EXPECT().GetByID(sEvent.ID).Return(sEvent, nil).Times(1)
		updatedEvent.Title = "Test"
		mockStorage.EXPECT().Edit(sEvent.ID, updatedEvent).Return(nil).Times(1)

		title := "Test"
		jsonEvent, errMarshal := json.Marshal(&EventRequest{ID: sEvent.ID, Title: &title})
		require.NoError(t, errMarshal)

		req := httptest.NewRequest(http.MethodPost, "/edit/", bytes.NewReader(jsonEvent))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("accept", "application/json")
		w := httptest.NewRecorder()
		s.edit(w, req)

		result := w.Result()
		defer result.Body.Close()

		buf := bufio.NewReader(result.Body)
		response, errBuf := buf.ReadString('\n')

		require.ErrorIs(t, errBuf, io.EOF)
		require.Equal(t, "", response)
		require.Equal(t, http.StatusOK, result.StatusCode)
	})

	t.Run("Check edit not found", func(t *testing.T) {
		mockStorage.EXPECT().GetByID(int64(0)).Return(storage.Event{}, storage.ErrEventNotFound).Times(1)

		title := "Edit title"
		jsonEvent, errMarshal := json.Marshal(&EventRequest{ID: 0, Title: &title})
		require.NoError(t, errMarshal)

		req := httptest.NewRequest(http.MethodPost, "/edit/", bytes.NewReader(jsonEvent))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("accept", "application/json")
		w := httptest.NewRecorder()
		s.edit(w, req)

		result := w.Result()
		defer result.Body.Close()

		response := &Response{}

		errDecode := json.NewDecoder(result.Body).Decode(response)
		require.NoError(t, errDecode)

		require.Contains(t, response.Error.Message, storage.ErrEventNotFound.Error())
		require.Equal(t, response.Data, nil)
		require.Equal(t, http.StatusBadRequest, result.StatusCode)
	})

	t.Run("Check wrong method type", func(t *testing.T) {
		jsonEvent, errMarshal := json.Marshal(&DeleteRequest{ID: 1})
		require.NoError(t, errMarshal)

		req := httptest.NewRequest(http.MethodGet, "/delete/", bytes.NewReader(jsonEvent))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("accept", "application/json")
		w := httptest.NewRecorder()
		s.edit(w, req)

		result := w.Result()
		defer result.Body.Close()

		response := &Response{}

		errDecode := json.NewDecoder(result.Body).Decode(response)
		require.NoError(t, errDecode)

		require.Contains(t, response.Error.Message, "method GET not not supported on uri /delete/")
		require.Equal(t, response.Data, nil)
		require.Equal(t, http.StatusMethodNotAllowed, result.StatusCode)
	})

	t.Run("Check delete", func(t *testing.T) {
		mockStorage.EXPECT().Delete(gomock.Any()).Return(nil).Times(1)

		jsonEvent, errMarshal := json.Marshal(&DeleteRequest{ID: 15})
		require.NoError(t, errMarshal)

		req := httptest.NewRequest(http.MethodPost, "/delete/", bytes.NewReader(jsonEvent))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("accept", "application/json")
		w := httptest.NewRecorder()
		s.delete(w, req)

		result := w.Result()
		defer result.Body.Close()

		buf := bufio.NewReader(result.Body)
		response, errBuf := buf.ReadString('\n')

		require.ErrorIs(t, errBuf, io.EOF)
		require.Equal(t, "", response)
		require.Equal(t, http.StatusOK, result.StatusCode)
	})

	t.Run("Check ListEventsDay", func(t *testing.T) {
		date := "2022-10-01"
		sEvent, err := prepareStorageEvent(date, 1)
		require.NoError(t, err)
		mockStorage.EXPECT().ListEventsDay(date).Return([]storage.Event{sEvent}, nil).Times(1)

		req := httptest.NewRequest(http.MethodGet, "/listeventsday/?date="+date, nil)
		req.Header.Set("accept", "application/json")
		w := httptest.NewRecorder()
		s.listEventsDay(w, req)

		result := w.Result()
		defer result.Body.Close()

		response := &Response{}
		response.Data = &ListResponse{}

		errDecode := json.NewDecoder(result.Body).Decode(response)
		require.NoError(t, errDecode)

		listResponse, ok := response.Data.(*ListResponse)
		if !ok {
			listResponse = &ListResponse{}
		}

		require.Equal(t, response.Error.Message, "")
		require.Equal(t, sEvent, s.convertRequestEvent(storage.Event{}, &listResponse.Events[0]))
		require.Equal(t, http.StatusOK, result.StatusCode)
	})

	t.Run("Check ListEventsWeek", func(t *testing.T) {
		date := "2022-10-01"
		date2 := "2022-10-06"
		sEvent1, err := prepareStorageEvent(date, 1)
		require.NoError(t, err)
		sEvent2, err := prepareStorageEvent(date2, 1)
		require.NoError(t, err)
		mockStorage.EXPECT().ListEventsWeek(date).Return([]storage.Event{sEvent1, sEvent2}, nil).Times(1)

		req := httptest.NewRequest(http.MethodGet, "/listeventsweek/?weekBeginDate="+date, nil)
		req.Header.Set("accept", "application/json")
		w := httptest.NewRecorder()
		s.listEventsWeek(w, req)

		result := w.Result()
		defer result.Body.Close()

		response := &Response{}
		response.Data = &ListResponse{}

		errDecode := json.NewDecoder(result.Body).Decode(response)
		require.NoError(t, errDecode)

		listResponse, ok := response.Data.(*ListResponse)
		if !ok {
			listResponse = &ListResponse{}
		}

		require.Equal(t, response.Error.Message, "")
		require.Equal(t, sEvent1, s.convertRequestEvent(storage.Event{}, &listResponse.Events[0]))
		require.Equal(t, sEvent2, s.convertRequestEvent(storage.Event{}, &listResponse.Events[1]))
		require.Equal(t, http.StatusOK, result.StatusCode)
	})

	t.Run("Check ListEventsMonth", func(t *testing.T) {
		date := "2022-10-01"
		date2 := "2022-10-30"
		sEvent1, err := prepareStorageEvent(date, 1)
		require.NoError(t, err)
		sEvent2, err := prepareStorageEvent(date2, 1)
		require.NoError(t, err)
		mockStorage.EXPECT().ListEventsMonth(date).Return([]storage.Event{sEvent1, sEvent2}, nil).Times(1)

		req := httptest.NewRequest(http.MethodGet, "/listeventsmonth/?monthBeginDate="+date, nil)
		req.Header.Set("accept", "application/json")
		w := httptest.NewRecorder()
		s.listEventsMonth(w, req)

		result := w.Result()
		defer result.Body.Close()

		response := &Response{}
		response.Data = &ListResponse{}

		errDecode := json.NewDecoder(result.Body).Decode(response)
		require.NoError(t, errDecode)

		listResponse, ok := response.Data.(*ListResponse)
		if !ok {
			listResponse = &ListResponse{}
		}

		require.Equal(t, response.Error.Message, "")
		require.Equal(t, sEvent1, s.convertRequestEvent(storage.Event{}, &listResponse.Events[0]))
		require.Equal(t, sEvent2, s.convertRequestEvent(storage.Event{}, &listResponse.Events[1]))
		require.Equal(t, http.StatusOK, result.StatusCode)
	})
}

func prepareEventRequest(addMinutes int) (EventRequest, error) {
	evRequest := EventRequest{}
	errFacker := gofake.Struct(&evRequest)

	*evRequest.DateStart = time.Now().Add(time.Minute * time.Duration(addMinutes)).Format(storage.DateTimeFormatISO)
	*evRequest.DateEnd = time.Now().Add(time.Minute * time.Duration(1+addMinutes)).Format(storage.DateTimeFormatISO)

	return evRequest, errFacker
}

func prepareStorageEvent(date string, num int) (storage.Event, error) {
	fEvent := storage.Event{}
	errFacker := gofake.Struct(&fEvent)
	if errFacker != nil {
		return storage.Event{}, errFacker
	}
	fEvent.Description = sql.NullString{
		String: gofake.Sentence(20),
		Valid:  true,
	}
	fEvent.NotificationInterval = sql.NullInt32{
		Int32: int32(gofake.IntRange(1, 4294967296)),
		Valid: true,
	}
	dt, err := time.ParseInLocation(storage.DateFormatISO, date, time.Local)
	if err != nil {
		return storage.Event{}, err
	}
	fEvent.DateStart = dt.Add(time.Minute * time.Duration(num))
	fEvent.DateEnd = dt.Add(time.Minute * time.Duration(num+1))

	return fEvent, errFacker
}
