package main

import (
	"github.com/stretchr/testify/require"
	"io/fs"
	"io/ioutil"
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	t.Run("wrong source type", func(t *testing.T) {
		err := Copy("/dev/urandom", "result/result.txt", 0, 0)
		require.Equal(t, ErrUnsupportedFile, err)
	})

	t.Run("wrong offset", func(t *testing.T) {
		err := Copy("testdata/input.txt", "result/result.txt", 7000, 0)
		require.Equal(t, ErrOffsetExceedsFileSize, err)
	})

	t.Run("file not exists", func(t *testing.T) {
		var pathError *fs.PathError

		err := Copy("testdata/randomfile.txt", "result/result.txt", 0, 0)
		require.IsType(t, pathError, err)
	})

	t.Run("multibyte copy", func(t *testing.T) {
		source := "testdata/input_multibyte.txt"
		compare := "testdata/input_multibyte_offset15_limit52.txt"
		dest := "result/result.txt"

		_ = Copy(source, dest, 15, 52)

		compareFile, _ := os.Open(compare)
		compareContent, _ := ioutil.ReadAll(compareFile)

		resultFile, _ := os.Open(dest)
		resultContent, _ := ioutil.ReadAll(resultFile)

		compareFile.Close()
		resultFile.Close()

		os.Remove(dest)

		require.Equal(t, compareContent, resultContent)
	})
}
