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

	client := NewTelnetClient(net.JoinHostPort(args[0], args[1]), timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		defer stop()
		err := client.Send()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
		}
	}()

	go func() {
		defer stop()
		err := client.Receive()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
		}
	}()

	<-ctx.Done()
}
