package memorystorage

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	gofake "github.com/brianvoe/gofakeit/v6"
	"github.com/director74/homeworks/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	ctx := context.Background()
	t.Run("Correct find event", func(t *testing.T) {
		st := New()

		fEvent, errPrepare := prepareItem(1)
		id, errAdd := st.Add(ctx, fEvent)

		require.NoError(t, errPrepare)
		require.NoError(t, errAdd)
		require.NotEqual(t, int64(0), id)

		foundedEvent, errSearch := checkExist(ctx, st, time.Now().Format(storage.DateFormatISO), id)

		require.NoError(t, errSearch)
		require.Equal(t, foundedEvent.ID, id)
	})

	t.Run("Correct delete event", func(t *testing.T) {
		var lastID int64
		st := New()

		for i := 1; i <= 2; i++ {
			fEvent, errPrepare := prepareItem(i)
			id, errAdd := st.Add(ctx, fEvent)

			lastID = id
			require.NotEqual(t, int64(0), id)
			require.NoError(t, errAdd)
			require.NoError(t, errPrepare)
		}

		err := st.Delete(ctx, lastID)
		require.NoError(t, err)
		require.Equal(t, 1, st.count())
	})

	t.Run("Correct edit event", func(t *testing.T) {
		var errEdit error
		st := New()

		fEvent, errPrepare := prepareItem(1)
		id, errAdd := st.Add(ctx, fEvent)

		require.NoError(t, errPrepare)
		require.NoError(t, errAdd)
		require.NotEqual(t, int64(0), id)

		foundedEvent, errGetList := checkExist(ctx, st, time.Now().Format(storage.DateFormatISO), id)

		require.NoError(t, errGetList)
		require.Equal(t, foundedEvent.ID, id)

		foundedEvent.Title = "New name for event"
		errEdit = st.Edit(ctx, id, foundedEvent)

		require.NoError(t, errEdit)

		foundedEvent, errGetList = checkExist(ctx, st, time.Now().Format(storage.DateFormatISO), id)

		require.NoError(t, errGetList)
		require.Equal(t, foundedEvent.ID, id)
		require.Equal(t, foundedEvent.Title, "New name for event")
	})

	t.Run("Error empty title", func(t *testing.T) {
		st := New()

		fEvent, errPrepare := prepareItem(1)
		fEvent.Title = ""
		id, errAdd := st.Add(ctx, fEvent)

		require.NoError(t, errPrepare)
		require.Equal(t, int64(0), id)
		require.ErrorIs(t, errAdd, storage.ErrWrongTitle)
	})

	t.Run("Error wrong datestart", func(t *testing.T) {
		st := New()

		fEvent, errPrepare := prepareItem(1)
		fEvent.DateStart = time.Now().Add(time.Duration(-5) * time.Hour)
		id, errAdd := st.Add(ctx, fEvent)

		require.NoError(t, errPrepare)
		require.Equal(t, int64(0), id)
		require.ErrorIs(t, errAdd, storage.ErrWrongDateStart)
	})

	t.Run("Error date busy", func(t *testing.T) {
		st := New()

		fEvent, errPrepare := prepareItem(1)
		fEvent.DateStart = time.Now().Add(time.Second * time.Duration(1))
		id, errAdd := st.Add(ctx, fEvent)

		require.NoError(t, errPrepare)
		require.NotEqual(t, int64(0), id)
		require.NoError(t, errAdd)

		fEvent, errPrepare = prepareItem(1)
		fEvent.DateStart = time.Now().Add(time.Second * time.Duration(1))
		id, errAdd = st.Add(ctx, fEvent)

		require.NoError(t, errPrepare)
		require.Equal(t, int64(0), id)
		require.ErrorIs(t, errAdd, storage.ErrDateBusy)
	})
}

func TestStorage_Add(t *testing.T) {
	st := New()
	ctx := context.Background()

	totalTries := 20

	t.Cleanup(func() {
		require.Equal(t, totalTries, st.count())
	})

	allEvents := make([]storage.Event, 0)
	for i := 1; i <= totalTries; i++ {
		fEvent, errPrepare := prepareItem(i)
		require.NoError(t, errPrepare)
		allEvents = append(allEvents, fEvent)
	}

	for i, event := range allEvents {
		event := event
		t.Run(fmt.Sprintf("Parallel add event %d", i), func(t *testing.T) {
			t.Parallel()
			id, errAdd := st.Add(ctx, event)
			require.NotEqual(t, int64(0), id)
			require.NoError(t, errAdd)
		})
	}
}

func prepareItem(num int) (storage.Event, error) {
	fEvent := storage.Event{}
	errFacker := gofake.Struct(&fEvent)
	fEvent.Description = sql.NullString{
		String: gofake.Sentence(20),
		Valid:  true,
	}
	fEvent.NotificationInterval = sql.NullInt32{
		Int32: int32(gofake.IntRange(1, 4294967296)),
		Valid: true,
	}
	fEvent.DateStart = time.Now().Add(time.Minute * time.Duration(num))
	fEvent.DateEnd = fEvent.DateStart.Add(time.Second * 2)

	return fEvent, errFacker
}

func checkExist(ctx context.Context, st *Storage, filterDate string, searchID int64) (storage.Event, error) {
	var foundedEvent storage.Event
	eventsSecond, err := st.ListEventsDay(ctx, filterDate)

	founded := false
	for _, rangeEvent := range eventsSecond {
		if rangeEvent.ID == searchID {
			foundedEvent = rangeEvent
			founded = true
		}
	}

	if !founded {
		return foundedEvent, fmt.Errorf("event not found")
	}

	return foundedEvent, err
}
