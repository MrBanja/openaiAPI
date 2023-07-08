package openaiAPI

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/MrBanja/openaiAPI/internal/field"
	inresponse "github.com/MrBanja/openaiAPI/internal/response"
	"net/http"
	"time"
)

type OpenAI struct {
	token   field.Password
	model   Model
	timeout time.Duration
}

func New(token string, model Model, timeout time.Duration) *OpenAI {
	return &OpenAI{
		token:   field.Password(token),
		model:   model,
		timeout: timeout,
	}
}

// SendWithStream sends a prompt with conversation history to OpenAI and streams responses back.
// Usage:
//
//	client := New("sk-xxx", Model4, 60*time.Second)
//	resp := client.SendWithStream(context.TODO(), "Tell me a joke", []Message{})
//
//	for {
//		select {
//		case err := <-resp.Error():
//			log.Panic(err)
//		case msg, ok := <-resp.Data():
//			if !ok {
//				return
//			}
//			fmt.Print(msg)
//		}
//	}
//
// Note that ctx cancellation's error will be sent to resp.Error() channel.
func (o *OpenAI) SendWithStream(ctx context.Context, prompt string, messages []Message) *ResponseStream {
	ctx, cancel := context.WithTimeout(ctx, o.timeout)

	responseStream := newResponseStream()
	messages = append(messages, Message{Role: RoleUser, Content: prompt})

	request := NewRequest(messages, true, o.model)
	data, err := json.Marshal(request)
	if err != nil {
		responseStream.sendError(err)
		cancel()
		return responseStream
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, urlV1Chat, bytes.NewReader(data))
	if err != nil {
		responseStream.sendError(err)
		cancel()
		return responseStream
	}

	req.Header.Set("Authorization", "Bearer "+o.token.Reveal())
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		responseStream.sendError(err)
		cancel()
		return responseStream
	}

	reader := bufio.NewReader(resp.Body)
	go func() {
		defer cancel()
		doneCh := make(chan struct{})
		go func() {
			defer resp.Body.Close()
			defer responseStream.close()
			defer close(doneCh)
			for {
				line, err := reader.ReadBytes('\n')
				if err != nil {
					responseStream.sendError(err)
					return
				}
				if !bytes.HasPrefix(line, []byte("data")) {
					continue
				}
				body := bytes.TrimPrefix(line, []byte("data: "))

				if string(body) == "[DONE]\n" {
					break
				}

				var stream inresponse.Stream
				if err := json.Unmarshal(body, &stream); err != nil {
					responseStream.sendError(err)
					return
				}

				msg := messageFromInternalResponseStream(stream)
				if msg == nil {
					responseStream.sendError(errors.New("empty response from OpenAI"))
					return
				}

				responseStream.send(msg.Content)
			}
		}()

		for {
			select {
			case <-ctx.Done():
				responseStream.sendError(ctx.Err())
				return
			case <-doneCh:
				return
			}
		}
	}()

	return responseStream
}
