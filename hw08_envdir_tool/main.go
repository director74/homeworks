package main

import "fmt"

func main() {
	//args := os.Args
	mp, err := ReadDir("testdata/env")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v", mp)

}
