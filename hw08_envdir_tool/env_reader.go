package main

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var ForbiddenFileSymbols = errors.New("filename contains =")

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	result := make(Environment, 0)

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		founded := strings.Index(file.Name(), "=")
		if founded >= 0 {
			return nil, ForbiddenFileSymbols
		}

		curFile, openErr := os.Open(dir + string(os.PathSeparator) + file.Name())
		if openErr != nil {
			return nil, openErr
		}

		sc := bufio.NewScanner(curFile)
		success := sc.Scan()
		if success {
			line := sc.Text()
			line = strings.TrimRight(line, "\t ")
			line = string(bytes.Replace([]byte(line), []byte("\000"), []byte("\n"), -1))
			result[file.Name()] = EnvValue{Value: line, NeedRemove: false}
		} else {
			result[file.Name()] = EnvValue{NeedRemove: true}
		}
		curFile.Close()
	}
	return result, nil
}
