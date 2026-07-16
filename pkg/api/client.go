package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	provider    string
	apiKey      string
	baseURL     string
	model       string
	timeout     time.Duration
	retryCount  int
	httpClient  *http.Client
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Temperature float64 `json:"temperature,omitempty"`
	MaxTokens    int     `json:"max_tokens,omitempty"`
	Stream     bool     `json:"stream,omitempty"`
}

type ChatResponse struct {
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Message Message `json:"message"`
}

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func New(provider, apiKey, baseURL, model string, timeout int, retryCount int) *Client {
	return &Client{
		provider:   provider,
		apiKey:     apiKey,
		baseURL:    baseURL,
		model:      model,
		timeout:    time.Duration(timeout) * time.Second,
		retryCount: retryCount,
		httpClient: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

func (c *Client) Chat(messages []Message) (string, error) {
	var lastErr error
	for attempt := 0; attempt <= c.retryCount; attempt++ {
		if attempt > 0 {
			delay := time.Duration(1<<uint(attempt-1)) * 2 * time.Second
			time.Sleep(delay)
		}

		result, err := c.doRequest(messages)
		if err == nil {
			return result, nil
		}
		lastErr = err
	}
	return "", lastErr
}

func (c *Client) doRequest(messages []Message) (string, error) {
	switch c.provider {
	case "anthropic":
		return c.anthropicChat(messages)
	default:
		return c.openaiChat(messages)
	}
}

func (c *Client) openaiChat(messages []Message) (string, error) {
	req := ChatRequest{
		Model:     c.model,
		Messages:  messages,
		Temperature: 0.0,
		MaxTokens: 8192,
		Stream:    false,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	url := c.baseURL
	if url[len(url)-1] != '/' {
		url += "/"
	}
	url += "chat/completions"

	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr APIError
		if err := json.Unmarshal(respBody, &apiErr); err == nil {
			return "", fmt.Errorf("API error (%d): %s", apiErr.Code, apiErr.Message)
		}
		return "", fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return "", err
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("empty response from API")
	}

	return chatResp.Choices[0].Message.Content, nil
}

func (c *Client) anthropicChat(messages []Message) (string, error) {
	var anthropicMessages []AnthropicMessage
	for _, m := range messages {
		anthropicMessages = append(anthropicMessages, AnthropicMessage{
			Role:    m.Role,
			Content: m.Content,
		})
	}

	req := AnthropicRequest{
		Model:       c.model,
		Messages:    anthropicMessages,
		System:      "You are a professional code auditor. Analyze code files thoroughly and provide structured reports.",
		Temperature: 0.0,
		MaxTokens:   8192,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	url := c.baseURL
	if url[len(url)-1] != '/' {
		url += "/"
	}
	url += "messages"

	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Anthropic API error (%d): %s", resp.StatusCode, string(respBody))
	}

	var anthropicResp AnthropicResponse
	if err := json.Unmarshal(respBody, &anthropicResp); err != nil {
		return "", err
	}

	if len(anthropicResp.Content) == 0 {
		return "", fmt.Errorf("empty response from Anthropic API")
	}

	return anthropicResp.Content[0].Text, nil
}
