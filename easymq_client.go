package goeasymq

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"net"
)

const (
	easymqPush          = 1
	easymqReadLatest    = 2
	easymqOk            = 1
	easymqInternalError = 2
)

type easymqClient struct {
	conn *net.TCPConn
}

func init() {
	clientInitFuncMap["easymq"] = func(s string) (Client, error) {
		addr, err := net.ResolveTCPAddr("tcp4", s)
		if err != nil {
			return nil, err
		}
		conn, err := net.DialTCP("tcp4", nil, addr)
		if err != nil {
			return nil, err
		}
		return &easymqClient{
			conn: conn,
		}, nil
	}
}
func (c *easymqClient) Push(ctx context.Context, topic string, content string) error {
	var topicLen, contentLen = len(topic), len(content)
	var buff = make([]byte, topicLen+contentLen+5)
	buff[0] = easymqPush
	binary.BigEndian.PutUint16(buff[1:3], uint16(topicLen))
	copy(buff[3:3+topicLen], []byte(topic))
	binary.BigEndian.PutUint16(buff[3+topicLen:5+topicLen], uint16(contentLen))
	copy(buff[5+topicLen:], []byte(content))
	_, err := c.conn.Write(buff)
	if err != nil {
		return err
	}
	var rdbuff [128]byte
	statusCode, size, err := readMsg(c.conn, rdbuff[:])
	if err != nil {
		return err
	}
	if statusCode != easymqOk {
		return errors.New(string(rdbuff[:size]))
	}
	return nil
}

func (c *easymqClient) ReadLatest(ctx context.Context, topic string) ([]byte, error) {
	var topicLen = len(topic)
	if topicLen == 0 {
		panic("topic is null")
	}
	var buff [1500]byte
	buff[0] = easymqReadLatest
	binary.BigEndian.PutUint16(buff[1:3], uint16(topicLen))
	copy(buff[3:topicLen+3], []byte(topic))
	_, err := c.conn.Write(buff[:topicLen+3])
	if err != nil {
		return nil, err
	}
	statusCode, size, err := readMsg(c.conn, buff[:])
	if err != nil {
		return nil, err
	}
	if statusCode != easymqOk {
		return nil, errors.New(string(buff[:size]))
	}
	return buff[:size], nil
}

func (c *easymqClient) Close() error {
	return c.conn.Close()
}
func readMsg(con io.Reader, buff []byte) (uint8, int, error) {
	var (
		readSize int
		size     int
		err      error
	)
readLen:
	size, err = con.Read(buff[:3-readSize])
	if err != nil || size == 0 {
		if err == nil {
			err = io.EOF
		}
		return 0, 0, err
	}
	readSize += size
	if readSize < 3 {
		goto readLen
	}
	statusCode := buff[0]
	usize := int(binary.BigEndian.Uint16(buff[1:3]))
	if usize == 0 {
		return statusCode, usize, nil
	}
	readSize = 0
readContent:
	size, err = con.Read(buff[:usize-readSize])
	if err != nil || size == 0 {
		if err == nil {
			err = io.EOF
		}
		return 0, 0, err
	}
	readSize += size
	if readSize < usize {
		goto readContent
	}
	return statusCode, usize, nil
}
