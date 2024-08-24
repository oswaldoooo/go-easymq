package goeasymq

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type httpClient struct {
	host   string
	client *http.Client
}

func init() {
	clientInitFuncMap["http"] = func(s string) (Client, error) {
		if s[len(s)-1] != '/' {
			s += "/"
		}
		return &httpClient{
			host:   s,
			client: &http.Client{},
		}, nil
	}
}
func (c *httpClient) Push(ctx context.Context, topic string, content string) error {
	rawcontent, _ := json.Marshal(messageReq{Topic: topic, Content: content})
	buff := bytes.NewReader(rawcontent)
	req, err := http.NewRequest(http.MethodPost, c.host+"publish", buff)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	raw, _ := io.ReadAll(resp.Body)
	return errors.New(string(raw))
}

func (c *httpClient) ReadLatest(ctx context.Context, topic string) ([]byte, error) {
	resp, err := c.client.Get(c.host + "read?topic=" + topic)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		return nil, errors.New(string(raw))
	}
	var objresp httpResponse
	err = json.NewDecoder(resp.Body).Decode(&objresp)
	if err != nil {
		return nil, errors.New("protocol version error " + err.Error())
	}
	if !objresp.Ok {
		return nil, errors.New(objresp.Reason)
	}
	return []byte(objresp.Data.Content), nil
}

func (c *httpClient) Close() error {
	return nil
}

type httpResponse struct {
	Ok     bool            `json:"ok,omitempty"`
	Reason string          `json:"reason,omitempty"`
	Data   messageResponse `json:"data,omitempty"`
}
type messageResponse struct {
	ID      uint32 `json:"id,omitempty"`
	Topic   string `json:"topic,omitempty"`
	Content string `json:"content,omitempty"`
}
