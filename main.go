package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

type Config struct {
	Mode             string
	RemoteAddr       string
	RemoteServerAddr string
	LocalAddr        string
}

var ServerConfig Config

func main() {

	config_file := flag.String("config", "", "config file")
	flag.Parse()

	if *config_file == "" {
		*config_file = "config.json"
	}

	f, e := os.OpenFile(*config_file, os.O_CREATE, os.ModePerm)
	if e != nil {
		log.Println("打开文件错误", *f)
		return
	}

	ServerConfig := Config{}
	err := json.NewDecoder(f).Decode(&ServerConfig)
	if err != nil {
		log.Println("解析配置文件错误", err)
		return
	}

	if ServerConfig.Mode == "client" {
		clientConfig := &ngClientConfig{ServerConfig.RemoteAddr, ServerConfig.LocalAddr}
		client := newNgClient(clientConfig)
		client.forever()
	} else {
		servConfig := &ngServerConfig{ServerConfig.RemoteServerAddr, ServerConfig.RemoteAddr}
		server := newNgServer(servConfig)
		server.start()
	}
}
