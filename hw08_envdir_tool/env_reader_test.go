package main

import (
	"errors"
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("wrong directory", func(t *testing.T) {
		envs, err := ReadDir("testdata/envv")

		require.Empty(t, envs)
		require.True(t, errors.Is(err, fs.ErrNotExist))
	})

	t.Run("empty directory path", func(t *testing.T) {
		envs, err := ReadDir("")

		require.Empty(t, envs)
		require.Equal(t, err, ErrEmptyDirPath)
	})

	t.Run("env name with =", func(t *testing.T) {
		dir, errMkDir := os.MkdirTemp("/tmp", "env_test")
		badFile, errTmpFile := os.CreateTemp(dir, "WORLD=")

		envs, err := ReadDir(dir)

		os.Remove(badFile.Name())
		errDirRemove := os.Remove(dir)

		require.Empty(t, envs)
		require.NoError(t, errMkDir)
		require.NoError(t, errTmpFile)
		require.NoError(t, errDirRemove)
		require.Equal(t, err, ErrForbiddenFileSymbols)
	})
}
