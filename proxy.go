package main

import (
	"log"
	"net"
	"sync"
)

type ProxyServer struct {
	servAddr  string
	proxyAddr string
}

func newProxyServer() *ProxyServer {
	return &ProxyServer{
		servAddr:  "127.0.0.1:18090",
		proxyAddr: "127.0.0.1:18080",
	}
}

func (p *ProxyServer) start() {
	l, err := net.Listen("tcp", p.servAddr)
	if err != nil {
		log.Panicln("监听端口错误", p.servAddr)
	}

	for {
		c, _ := l.Accept()
		p.servConn(c)
	}
}

func (p *ProxyServer) servConn(c net.Conn) {
	defer c.Close()

	lc, err := net.Dial("tcp", p.proxyAddr)
	if err != nil {
		log.Panicln("链接到远程失败", err, p.proxyAddr)
	}

	defer lc.Close()

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		buff := make([]byte, 8192)

		for {
			n, err := lc.Read(buff)
			if err != nil {
				break
			}
			_, err = c.Write(buff[:n])
			if err != nil {
				break
			}
		}
	}()

	go func() {
		defer wg.Done()
		buff := make([]byte, 8192)

		for {
			n, err := c.Read(buff)

			if err != nil {
				break
			}
			_, err = lc.Write(buff[:n])
			if err != nil {
				break
			}
		}
	}()

	wg.Wait()
}
