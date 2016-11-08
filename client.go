package main

import (
	"bufio"
	"net"
	"sync"
)

type ngClientConfig struct {
	servAddr   string
	clientAddr string
}

type ngClient struct {
	servs  map[int64]*ngConn
	sMutex sync.Mutex

	client     net.Conn // 客户端链接
	servAddr   string
	clientAddr string
}

func newNgClient(config *ngClientConfig) *ngClient {
	server := &ngClient{
		servs:      make(map[int64]*ngConn),
		sMutex:     sync.Mutex{},
		client:     nil,
		clientAddr: config.clientAddr,
		servAddr:   config.servAddr,
	}
	return server
}

func (s *ngClient) start() {
	c, err := net.Dial("tcp", s.servAddr)
	if err != nil {
		return
	}

	defer c.Close()

	s.client = c
	reader := bufio.NewReader(c)

	for {
		msg, err := readMsg(reader)

		if err != nil {
			break
		}

		switch msg.head {
		case HEAD_NEWCONN:
			seq := msg.seq

			// 链接到本地服务器
			c, err := net.Dial("tcp", s.clientAddr)
			if err != nil {
				s.closeSeq(seq)
				continue
			}
			s.addServ(seq, c)

			go s.servLocal(seq)
		case HEAD_CLOSE:
			seq := msg.seq
			s.closeSeq(seq)
		case HEAD_CONTENT:
			seq := msg.seq
			r := s.getServ(seq)
			if r != nil {
				n, err := r.c.Write(msg.buff)
				if err != nil {
					s.closeSeq(seq)
					continue
				}
				if n != len(msg.buff) {
					s.closeSeq(seq)
					continue
				}
			}
		}
	}
}

func (s *ngClient) servLocal(seq int64) {
	c := s.getServ(seq)
	if c == nil {
		return
	}

	defer s.closeSeq(seq)

	if tc, ok := c.c.(*net.TCPConn); ok {
		tc.SetNoDelay(true)
	}

	buff := make([]byte, 1400)
	for {
		n, err := c.c.Read(buff)

		if err != nil {
			break
		}

		msg := newMsg(HEAD_CONTENT, seq, buff[:n])

		_, err = s.client.Write(encodeMsg(msg))
		if err != nil {
			break
		}
	}
}

func (s *ngClient) closeSeq(seq int64) {
	s.sMutex.Lock()
	defer s.sMutex.Unlock()

	if c, ok := s.servs[seq]; ok {
		delete(s.servs, seq)
		c.c.Close()

		closeMsg := newMsg(HEAD_CLOSE, seq, []byte{})
		s.client.Write(encodeMsg(closeMsg))
	}
}

func (s *ngClient) getServ(seq int64) *ngConn {
	s.sMutex.Lock()
	defer s.sMutex.Unlock()

	if c, ok := s.servs[seq]; ok {
		return c
	}
	return nil
}

func (s *ngClient) addServ(seq int64, c net.Conn) {
	s.sMutex.Lock()
	defer s.sMutex.Unlock()
	s.servs[seq] = &ngConn{c}
}
