package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type HttpServer struct{}

func (s *HttpServer) start() {
	err := http.ListenAndServe("127.0.0.1:8090", s)
	if err != nil {
		log.Println("HttpServer错误", err)
	}
}

func (s *HttpServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// 移除host
	r.Header.Del("host")

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: s,
	}

	r.RequestURI = ""
	resp, err := client.Do(r)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	rw.Write(buff)
}

func (s *HttpServer) RoundTrip(r *http.Request) (*http.Response, error) {
	// handle request

	return nil, io.EOF
}
