package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Client struct {
	baseURL string
	client  *http.Client
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Response struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func NewClient() *Client {
	baseURL := os.Getenv("LLM_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:11434/api/chat"
	}

	return &Client{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 2 * time.Minute,
		},
	}
}

func (c *Client) Query(prompt string) (string, error) {
	req := Request{
		Model: os.Getenv("LLM_MODEL"),
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	if req.Model == "" {
		req.Model = "llama3.2"
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.client.Post(c.baseURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return result.Choices[0].Message.Content, nil
}

func (c *Client) StreamQuery(prompt string, callback func(string)) error {
	req := Request{
		Model: os.Getenv("LLM_MODEL"),
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	if req.Model == "" {
		req.Model = "llama3.2"
	}

	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := c.client.Post(c.baseURL+"?stream=true", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	for {
		var chunk map[string]interface{}
		if err := decoder.Decode(&chunk); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to decode chunk: %w", err)
		}

		if message, ok := chunk["message"].(map[string]interface{}); ok {
			if content, ok := message["content"].(string); ok {
				callback(content)
			}
		}

		if done, ok := chunk["done"].(bool); ok && done {
			break
		}
	}

	return nil
}
