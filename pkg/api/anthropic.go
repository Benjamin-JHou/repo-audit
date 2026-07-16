package api

type AnthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AnthropicRequest struct {
	Model       string              `json:"model"`
	Messages    []AnthropicMessage  `json:"messages"`
	System      string              `json:"system"`
	Temperature float64             `json:"temperature,omitempty"`
	MaxTokens     int               `json:"max_tokens,omitempty"`
}

type AnthropicResponse struct {
	Content []AnthropicContent `json:"content"`
}

type AnthropicContent struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}
