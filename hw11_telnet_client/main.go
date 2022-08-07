package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"
)

var timeout time.Duration

var ErrNotEnoughArguments = errors.New("not enough arguments passed")

func init() {
	flag.DurationVar(&timeout, "timeout", time.Second*10, "timeout duration")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		fmt.Println(ErrNotEnoughArguments)
		return
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	client := NewTelnetClient(net.JoinHostPort(args[0], args[1]), timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	go func() {
		for {
			err := client.Send()
			if err != nil {
				fmt.Println(err)
				break
			}
		}
		stop()
	}()

	go func() {
		for {
			err := client.Receive()
			if err != nil {
				fmt.Println(err)
				break
			}
		}
		stop()
	}()

	<-ctx.Done()
}
