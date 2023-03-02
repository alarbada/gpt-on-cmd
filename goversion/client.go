package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCompletionReq struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Client struct {
	APIKey string
	client *http.Client
}

func NewClient(apikey string) *Client {
	return &Client{
		APIKey: apikey,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (this *Client) path(p string) string {
	return fmt.Sprintf("https://api.openai.com%s", p)
}

type ChatCompletionChunk struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Delta struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"delta"`
		Index        int    `json:"index"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

func (this *Client) ChatCompletion(completionReq ChatCompletionReq, onData func(string, error)) {
	data, err := json.Marshal(completionReq)
	if err != nil {
		onData("", err)
		return
	}

	body := bytes.NewBuffer(data)

	req, err := http.NewRequest("POST", this.path("/v1/chat/completions"), body)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+this.APIKey)

	res, err := this.client.Do(req)
	if err != nil {
		onData("", err)
		return
	}

	doneC := make(chan struct{})
	go func() {
		defer close(doneC)
		reader := bufio.NewReader(res.Body)
		for {
			line, err := reader.ReadString('\n')
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				onData("", err)
				return
			}

			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			// Trim the "data: " prefix
			line = line[6:]

			// if line includes "[DONE]" then we're done
			if strings.Contains(line, "[DONE]") {
				break
			}

			var chunk ChatCompletionChunk
			err = json.Unmarshal([]byte(line), &chunk)
			if err != nil {
				onData("", err)
				return
			}

			onData(chunk.Choices[0].Delta.Content, err)
		}
	}()

	<-doneC
}
