package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func __Test_msg(t *testing.T) {
	buff := []byte("hahaex")
	msg := newMsg(HEAD_CONTENT, 20, buff)
	logmsg(msg)
	enBuf := encodeMsg(msg)

	enBuf = append(enBuf, enBuf...)
	enBuf = append(enBuf, enBuf...)
	enBuf = append(enBuf, enBuf...)

	logmsg(enBuf)

	r := bufio.NewReader(bytes.NewReader(enBuf))
	msg2, _ := readMsg(r)
	msg3, _ := readMsg(r)
	msg4, _ := readMsg(r)
	logmsg(msg2.head, string(msg2.buff), msg2.length, msg2.seq)
	logmsg(msg3.head, string(msg3.buff), msg3.length, msg3.seq)
	logmsg(msg4.head, string(msg4.buff), msg2.length, msg4.seq)
}

func __Test_http(t *testing.T) {
	//resp, err := http.Get("")
	req, err := http.NewRequest("GET", "http://127.0.0.1:80/", nil)
	req.Header.Add("HOST", "127.0.0.1:8080")
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(body))
}

func Test_buff(t *testing.T) {
	str := "71 69 84 32 47 32 72 84 84 80 47 49 46 49 13 10 72 111 115 116 58 32 49 50 55 46 48 46 48 46 49 58 49 56 48 57 48 13 10 67 111 110 110 101 99 116 105 111 110 58 32 107 101 101 112 45 97 108 105 118 101 13 10 67 97 99 104 101 45 67 111 110 116 114 111 108 58 32 109 97 120 45 97 103 101 61 48 13 10 85 112 103 114 97 100 101 45 73 110 115 101 99 117 114 101 45 82 101 113 117 101 115 116 115 58 32 49 13 10 85 115 101 114 45 65 103 101 110 116 58 32 77 111 122 105 108 108 97 47 53 46 48 32 40 87 105 110 100 111 119 115 32 78 84 32 54 46 51 59 32 87 105 110 54 52 59 32 120 54 52 41 32 65 112 112 108 101 87 101 98 75 105 116 47 53 51 55 46 51 54 32 40 75 72 84 77 76 44 32 108 105 107 101 32 71 101 99 107 111 41 32 67 104 114 111 109 101 47 53 51 46 48 46 50 55 56 53 46 49 52 51 32 83 97 102 97 114 105 47 53 51 55 46 51 54 13 10 65 99 99 101 112 116 58 32 116 101 120 116 47 104 116 109 108 44 97 112 112 108 105 99 97 116 105 111 110 47 120 104 116 109 108 43 120 109 108 44 97 112 112 108 105 99 97 116 105 111 110 47 120 109 108 59 113 61 48 46 57 44 105 109 97 103 101 47 119 101 98 112 44 42 47 42 59 113 61 48 46 56 13 10 65 99 99 101 112 116 45 69 110 99 111 100 105 110 103 58 32 103 122 105 112 44 32 100 101 102 108 97 116 101 44 32 115 100 99 104 13 10 65 99 99 101 112 116 45 76 97 110 103 117 97 103 101 58 32 122 104 45 67 78 44 122 104 59 113 61 48 46 56 44 101 110 45 85 83 59 113 61 48 46 54 44 101 110 59 113 61 48 46 52 13 10 13 10"
	bytesStr := strings.Split(str, " ")
	buff := make([]byte, len(bytesStr))

	for i := 0; i < len(bytesStr); i++ {
		j, _ := strconv.ParseInt(bytesStr[i], 10, 8)
		buff[i] = byte(j)
	}

	log.Println(buff)

	c, _ := net.Dial("tcp", "127.0.0.1:80")
	c.Write(buff)

	buff2, _ := ioutil.ReadAll(c)
	log.Println(string(buff2))
}
