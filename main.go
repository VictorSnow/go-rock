package main

import (
	"time"
)

func main() {
	// 测试
	servConfig := &ngServerConfig{"127.0.0.1:18080", "127.0.0.1:18081"}
	clientConfig := &ngClientConfig{"127.0.0.1:18081", "127.0.0.1:80"}

	server := newNgServer(servConfig)
	client := newNgClient(clientConfig)

	go newProxyServer().start()

	server.start()
	client.start()

	for {
		time.Sleep(60 * time.Second)
	}
}
