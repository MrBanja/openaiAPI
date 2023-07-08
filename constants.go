package openaiAPI

const (
	RoleUser      Role = "user"
	RoleSystem    Role = "system"
	RoleAssistant Role = "assistant"
)

const (
	urlV1Chat = "https://api.openai.com/v1/chat/completions"
)

const (
	Model3 Model = "gpt-3.5-turbo"
	Model4 Model = "gpt-4"
)
