# OpenAI API

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/MrBanja/openaiAPI/blob/main/LICENSE)

This is a Go package for interacting with the OpenAI Stream API for chat completions. It allows you to easily integrate chat completion functionality into your Go applications.

## Installation

To use this package, you need to have Go installed and set up on your machine. Then, you can install the package using the following command:

```shell
go get github.com/MrBanja/openaiAPI
```

## Usage

### Getting Started
Here is a simplified example of how to use the OpenAI API package:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/MrBanja/openaiAPI"
)

func main() {
	// Create a new client
	client := openaiAPI.New("sk-xxx", openaiAPI.Model4, 60*time.Second)

	// Send a chat message and receive responses
	resp := client.SendWithStream(context.TODO(), "Tell me a joke", []openaiAPI.Message{})

	for {
		select {
		case err := <-resp.Error():
			// Error received or context cancelled
			log.Panic(err)
		case msg, ok := <-resp.Data():
			if !ok {
				// Stream closed. No more messages to receive
				return
			}
			fmt.Print(msg)
		}
	}
}
```

In this example, we first create a new client by providing your OpenAI Secret Key (`sk-xxx` in this case), the desired model (`Model4` in this case), and the timeout duration for the API requests (60 seconds in this case).

We then use the client to send a chat message using the `SendWithStream` method, passing the context, the message text ("Tell me a joke" in this case), and an empty list of previous messages.

Finally, we use a `for` loop with a `select` statement to handle incoming data or errors from the response. Any error received will cause the program to panic, and each received message will be printed to the console.

### With Previous Messages
Here is an example of how to use the OpenAI API package with previous messages:

```go
var messages []openaiAPI.Message
messages = append(
    messages,
    openaiAPI.Message{Role: openaiAPI.RoleUser, Content: "Hi"},
    openaiAPI.Message{Role: openaiAPI.RoleAssistant, Content: "Hello, how are you?"},
)
resp := client.SendWithStream(context.TODO(), "Tell me a joke", messages)
```

### Some real world examples
Here are some [real world](https://github.com/MrBanja/gcli/blob/b8cd6d7f49fbbcb89252254da3ab3b41713421c5/internal/controller/flow_controller.go#L25) examples of how to use the OpenAI API package:

```go
func (f *FlowController) Stream(prompt string) error {
	convID := viper.GetString("current_conversation_id")
	messages, err := f.history.Get(convID)
	if err != nil {
		return err
	}

	resp := f.openai.SendWithStream(context.Background(), prompt, messages)
	messageContent := ""

L:
	for {
		select {
		case err := <-resp.Error():
			util.HandleError(err, "OpenAI response error")
		case msg, ok := <-resp.Data():
			if !ok {
				break L
			}
			fmt.Print(msg)
			messageContent += msg
		}
	}

	screen.Clear()
	screen.MoveTopLeft()
	util.PrintOrExit("# RESPONSE")
	util.PrintOrExit(messageContent)

	messages = append(
		messages,
		openaiAPI.Message{Role: openaiAPI.RoleUser, Content: prompt},
		openaiAPI.Message{Role: openaiAPI.RoleAssistant, Content: messageContent},
	)
	if err := f.history.Set(convID, messages); err != nil {
		return nil
	}
	return nil
}
```

## License

This project is licensed under the [MIT license](https://github.com/MrBanja/openaiAPI/blob/main/LICENSE).
