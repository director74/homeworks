package main

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var (
	ErrForbiddenFileSymbols = errors.New("filename contains =")
	ErrEmptyDirPath         = errors.New("path to env files is empty")
)

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	result := make(Environment)

	if len(dir) == 0 {
		return nil, ErrEmptyDirPath
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		founded := strings.Contains(file.Name(), "=")
		if founded {
			return nil, ErrForbiddenFileSymbols
		}

		curFile, openErr := os.Open(filepath.Join(dir, string(os.PathSeparator), file.Name()))
		if openErr != nil {
			return nil, openErr
		}

		sc := bufio.NewScanner(curFile)
		success := sc.Scan()
		err := sc.Err()
		if err != nil {
			return nil, err
		}
		if success {
			line := sc.Text()
			line = strings.TrimRight(line, "\t ")
			line = string(bytes.ReplaceAll([]byte(line), []byte("\000"), []byte("\n")))
			result[file.Name()] = EnvValue{Value: line, NeedRemove: false}
		} else {
			result[file.Name()] = EnvValue{NeedRemove: true}
		}
		curFile.Close()
	}
	return result, nil
}
