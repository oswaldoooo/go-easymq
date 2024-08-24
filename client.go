package goeasymq

import (
	"context"
	"errors"
)

var clientInitFuncMap = make(map[string]func(string) (Client, error))

type Client interface {
	Push(ctx context.Context, topic string, content string) error
	ReadLatest(ctx context.Context, topic string) ([]byte, error)
	Close() error
}

func Connect(network string, address string) (Client, error) {
	initFunc, ok := clientInitFuncMap[network]
	if !ok {
		return nil, errors.New("not found network client " + network)
	}
	return initFunc(address)
}

type messageReq struct {
	Topic   string `json:"topic,omitempty"`
	Content string `json:"content,omitempty"`
}
