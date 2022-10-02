package logger

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	const logFormat = "2006/01/02 15:04:05"

	t.Run("Check warn message", func(t *testing.T) {
		bt := make([]byte, 0)
		buff := bytes.NewBuffer(bt)

		logg := New(logWarn)
		logg.SetOutput(buff)

		logg.Warn("Very important notification")

		str, err := buff.ReadString('\n')

		require.NoError(t, err)
		require.Equal(t, time.Now().Format(logFormat)+" log.WARN Very important notification\n", str)
	})

	t.Run("Check log level restriction", func(t *testing.T) {
		bt := make([]byte, 0)
		buff := bytes.NewBuffer(bt)

		logg := New(logError)
		logg.SetOutput(buff)

		logg.Info("Message that won't be shown")

		str, err := buff.ReadString('\n')
		require.Equal(t, "", str)
		require.ErrorIs(t, err, io.EOF)
	})
}
