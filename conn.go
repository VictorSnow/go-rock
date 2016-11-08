package main

import (
	"errors"
	"io"
	"net"
)

const NGCONN_RETRY = 5

var emptyIoError = errors.New("Error Empty io")

type emptyNgConn struct{}

func (r *emptyNgConn) Read(buff []byte) (int, error) {
	return 0, emptyIoError
}

func (r *emptyNgConn) Close() error {
	return emptyIoError
}

func (r *emptyNgConn) Write(buff []byte) (int, error) {
	return 0, emptyIoError
}

type ngConn struct {
	c net.Conn
}

func (c *ngConn) Read(buff []byte) (int, error) {
	retry := NGCONN_RETRY
	for retry != 0 {
		n, err := c.c.Read(buff)
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				retry--
				continue
			} else {
				return n, err
			}
		}
		return n, err
	}
	return 0, io.EOF
}

func (c *ngConn) Write(buff []byte) (int, error) {
	retry := NGCONN_RETRY
	for retry != 0 {
		n, err := c.c.Write(buff)
		if err != nil {
			if e, ok := err.(net.Error); ok && e.Timeout() {
				retry--
				continue
			} else {
				return n, err
			}
		}
		return n, err
	}
	return 0, io.EOF
}

func (c *ngConn) Close() error {
	return c.c.Close()
}
