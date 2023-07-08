package openaiAPI

type Request struct {
	Model    Model        `json:"model"`
	Messages []ReqMessage `json:"messages"`
	Stream   bool         `json:"stream"`
}

type ReqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewRequest(messages []Message, isStream bool, model Model) Request {
	return Request{
		Model: model,
		Messages: func() []ReqMessage {
			var reqMessages []ReqMessage
			for _, message := range messages {
				reqMessages = append(reqMessages, ReqMessage{
					Role:    string(message.Role),
					Content: message.Content,
				})
			}
			return reqMessages
		}(),
		Stream: isStream,
	}
}
