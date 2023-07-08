package openaiAPI

import inresponse "github.com/MrBanja/openaiAPI/internal/response"

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

func messageFromInternalResponseStream(s inresponse.Stream) *Message {
	if len(s.Choices) == 0 {
		return nil
	}
	return &Message{
		Role:    Role(s.Choices[0].Delta.Role),
		Content: s.Choices[0].Delta.Content,
	}
}
