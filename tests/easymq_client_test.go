package tests

import (
	"context"
	"testing"

	goeasymq "github.com/oswaldoooo/go-easymq"
)

var easymqClient goeasymq.Client

func init() {
	var err error
	easymqClient, err = goeasymq.Connect("easymq", "localhost:7777")
	if err != nil {
		panic(err)
	}
}
func TestEasyMqPush(t *testing.T) {
	err := easymqClient.Push(context.Background(), "testGo", "hello easymq")
	if err != nil {
		t.Fatal(err)
	}
}
func TestEasyMqReadLatest(t *testing.T) {
	content, err := easymqClient.ReadLatest(context.Background(), "testGo")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(content))
}
