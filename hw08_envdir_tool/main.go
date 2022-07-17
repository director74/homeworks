package main

import (
	"fmt"
	"os"
)

type args struct {
	dirPath string
	arg1    string
	arg2    string
}

func processArgs(params []string) args {
	result := args{}

	switch len(params) {
	case 2:
		result.dirPath = params[1]
	case 3:
		result.arg1 = params[2]
	case 4:
		result.arg2 = params[3]
	}

	return result
}

func main() {
	args := processArgs(os.Args)

	_, err := ReadDir(args.dirPath)
	if err != nil {
		fmt.Println(err)
	}
}
