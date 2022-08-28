package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	conn        net.Conn
	in          io.ReadCloser
	out         io.Writer
	connAddress string
	connTimeout time.Duration
}

func (t *telnetClient) Connect() error {
	var err error
	t.conn, err = net.DialTimeout("tcp", t.connAddress, t.connTimeout)
	if err != nil {
		return fmt.Errorf("cannot connect: %w", err)
	}
	return nil
}

func (t *telnetClient) Send() error {
	_, err := io.Copy(t.conn, t.in)
	if err != nil {
		return fmt.Errorf("cannot send request: %w", err)
	}

	if _, err = os.Stderr.Write([]byte("...EOF\n")); err != nil {
		return err
	}

	return nil
}

func (t *telnetClient) Receive() error {
	_, err := io.Copy(t.out, t.conn)
	if err != nil {
		return fmt.Errorf("cannot read response: %w", err)
	}

	if _, err = os.Stderr.Write([]byte("...Connection was closed by peer\n")); err != nil {
		return err
	}

	return nil
}

func (t *telnetClient) Close() error {
	if err := t.conn.Close(); err != nil {
		return fmt.Errorf("cannot close connection: %w", err)
	}
	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		connAddress: address,
		connTimeout: timeout,
		in:          in,
		out:         out,
	}
}
