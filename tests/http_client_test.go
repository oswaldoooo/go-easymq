package tests

import (
	"context"
	"testing"

	goeasymq "github.com/oswaldoooo/go-easymq"
)

var httpClient goeasymq.Client

func init() {
	var err error
	httpClient, err = goeasymq.Connect("http", "http://localhost:8080/")
	if err != nil {
		panic(err)
	}
}
func TestHttpReadLatest(t *testing.T) {
	content, err := httpClient.ReadLatest(context.Background(), "testGo")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(content))
}
func TestHttpPush(t *testing.T) {
	err := httpClient.Push(context.Background(), "testGo", "hello easymq")
	if err != nil {
		t.Fatal(err)
	}
}
