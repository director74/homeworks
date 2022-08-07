package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	conn           net.Conn
	in             io.ReadCloser
	out            io.Writer
	inBuf          []byte
	conBuf         []byte
	connAddress    string
	connTimeout    time.Duration
	responseReader *bufio.Reader
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
	nRead, err := t.in.Read(t.inBuf)

	if errors.Is(err, io.EOF) {
		return fmt.Errorf("...EOF")
	} else if err != nil {
		return fmt.Errorf("cannot process input: %w", err)
	}

	if nRead > 0 {
		_, err := t.conn.Write(t.inBuf[:nRead])
		if err != nil {
			return fmt.Errorf("cannot send request: %w", err)
		}
	}

	return nil
}

func (t *telnetClient) Receive() error {
	if t.responseReader == nil {
		t.responseReader = bufio.NewReader(t.conn)
	}
	nRead, err := t.responseReader.Read(t.conBuf)
	if errors.Is(err, io.EOF) {
		return fmt.Errorf("...Connection was closed by peer")
	} else if err != nil {
		return fmt.Errorf("cannot read response: %w", err)
	}

	if nRead > 0 {
		_, err := t.out.Write(t.conBuf[:nRead])
		if err != nil {
			return fmt.Errorf("cannot show response: %w", err)
		}
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
		inBuf:       make([]byte, 512),
		conBuf:      make([]byte, 512),
	}
}
