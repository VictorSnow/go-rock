package main

import (
	"bufio"
	"errors"
	"io"
)

const MAX_PACK_LEN = 1400

const HEAD_NEWCONN = 1
const HEAD_CONTENT = 2
const HEAD_CLOSE = 3
const HEAD_KEEPALIVE = 4

type ngMessage struct {
	head   byte
	length int16
	seq    int64
	buff   []byte
}

func newMsg(head byte, seq int64, buff []byte) (msg *ngMessage) {
	msg = &ngMessage{
		head:   head,
		seq:    seq,
		length: int16(len(buff)),
		buff:   buff,
	}
	return
}

func encodeMsg(msg *ngMessage) []byte {
	length := 11 + msg.length
	buff := make([]byte, length)
	buff[0] = msg.head
	buff[1] = byte(msg.length & 255)
	buff[2] = byte(msg.length >> 8)
	buff[3] = byte(msg.seq & 255)
	buff[4] = byte(msg.seq >> 8)
	buff[5] = byte(msg.seq >> 16)
	buff[6] = byte(msg.seq >> 24)
	buff[7] = byte(msg.seq >> 32)
	buff[8] = byte(msg.seq >> 40)
	buff[9] = byte(msg.seq >> 48)
	buff[10] = byte(msg.seq >> 56)

	copy(buff[11:], msg.buff)
	return buff
}

// 解析得到msg信息
func readMsg(r *bufio.Reader) (*ngMessage, error) {
	head := make([]byte, 11)
	err := safeRead(r, head)

	if err != nil {
		return nil, err
	}

	msg := &ngMessage{}
	msg.head = head[0]
	msg.length = int16(head[1]) + int16(head[2])<<8
	msg.seq = int64(head[3]) + int64(head[4])<<8 + int64(head[5])<<16 + int64(head[6])<<24 + int64(head[7])<<32 + int64(head[8])<<40 + int64(head[9])<<48 + int64(head[10])<<56
	msg.buff = []byte{}

	if msg.head != 1 && msg.head != 2 && msg.head != 3 {
		return nil, errors.New("消息头错误")
	}

	if msg.length > 0 {
		msg.buff = make([]byte, msg.length)
		err := safeRead(r, msg.buff)
		if err != nil {
			return nil, err
		}
	}
	return msg, nil
}

// 安全读取制定字节
func safeRead(r io.Reader, buff []byte) error {
	offset := 0
	length := len(buff)
	for {
		n, err := r.Read(buff[offset:])
		if err != nil {
			return err
		}
		offset += n
		if offset == length {
			return nil
		}
	}
}
