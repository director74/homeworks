package main

import (
	"errors"
	"github.com/stretchr/testify/require"
	"io/fs"
	"os"
	"testing"
)

func TestReadDir(t *testing.T) {
	t.Run("wrong directory", func(t *testing.T) {
		_, err := ReadDir("testdata/envv")

		require.True(t, errors.Is(err, fs.ErrNotExist))
	})

	t.Run("empty directory path", func(t *testing.T) {
		_, err := ReadDir("")

		require.Equal(t, err, ErrEmptyDirPath)
	})

	t.Run("env name with =", func(t *testing.T) {
		dir, _ := os.MkdirTemp("/tmp", "env_test")
		badFile, errTmpFile := os.CreateTemp(dir, "WORLD=")

		_, err := ReadDir(dir)

		os.Remove(badFile.Name())
		errDirRemove := os.Remove(dir)

		require.NoError(t, errTmpFile)
		require.NoError(t, errDirRemove)
		require.Equal(t, err, ErrForbiddenFileSymbols)
	})
}
