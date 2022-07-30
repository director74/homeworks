package main

import (
	"errors"
	"fmt"
	"os"
)

var ErrNotEnoughArguments = errors.New("not enough few arguments passed")

func main() {
	var dirPath string

	if len(os.Args) < 2 {
		fmt.Println(ErrNotEnoughArguments)
		return
	}
	dirPath = os.Args[1]

	envs, err := ReadDir(dirPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	RunCmd(os.Args[2:], envs)
}
