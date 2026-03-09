// Package ailib - ai-lib-go 主入口
// Official Go Runtime for AI-Protocol

package ailib

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hiddenpath/ai-lib-go/internal/protocol"
)

// Version is the library version.
const Version = "0.5.0"

// client implements the Client interface.
type client struct {
	manifest interface{} // v1.ProviderManifest or v2.ProviderManifest
	loader   *protocol.Loader
	http     *http.Client
	apiKey   string
	baseURL  string
	headers  map[string]string
}

// newClient creates a new client.
func newClient(b *ClientBuilder) (Client, error) {
	c := &client{
		loader: protocol.NewLoader(),
		http: &http.Client{
			Timeout: time.Duration(b.timeout) * time.Second,
		},
		apiKey:  b.apiKey,
		headers: b.headers,
	}

	// Load protocol
	if b.protocolData != nil {
		manifest, err := c.loader.LoadFromData(b.protocolData)
		if err != nil {
			return nil, fmt.Errorf("failed to load protocol data: %w", err)
		}
		c.manifest = manifest
	}

	// Set base URL
	if b.baseURL != "" {
		c.baseURL = b.baseURL
	} else if c.manifest != nil {
		endpoint, err := protocol.GetEndpoint(c.manifest)
		if err != nil {
			return nil, err
		}
		c.baseURL = endpoint
	}

	return c, nil
}

// Chat performs a synchronous chat completion.
func (c *client) Chat(ctx context.Context, messages []Message, opts *ChatOptions) (*ChatResponse, error) {
	if opts == nil {
		opts = &ChatOptions{}
	}

	// Build request
	reqBody := c.buildRequest(messages, opts)
	reqBody["stream"] = false

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Build HTTP request
	url := c.baseURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	c.setHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check for errors
	if resp.StatusCode >= 400 {
		return nil, c.parseError(resp)
	}

	// Parse response
	var response ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// ChatStream performs a streaming chat completion.
func (c *client) ChatStream(ctx context.Context, messages []Message, opts *ChatOptions) (Stream, error) {
	if opts == nil {
		opts = &ChatOptions{}
	}

	// Build request
	reqBody := c.buildRequest(messages, opts)
	reqBody["stream"] = true

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Build HTTP request
	url := c.baseURL + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	c.setHeaders(req)
	req.Header.Set("Content-Type", "application/json"

)
	// Send request
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Check for errors
	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		return nil, c.parseError(resp)
	}

	// Create stream reader
	return newSSEStream(resp.Body), nil
}

// Close releases resources.
func (c *client) Close() error {
	c.http.CloseIdleConnections()
	return nil
}

// buildRequest builds the request body.
func (c *client) buildRequest(messages []Message, opts *ChatOptions) map[string]interface{} {
	req := make(map[string]interface{})

	// Model
	if opts.Model != "" {
		req["model"] = opts.Model
	}

	// Messages
	req["messages"] = c.convertMessages(messages)

	// Options
	if opts.Temperature != nil {
		req["temperature"] = *opts.Temperature
	}
	if opts.MaxTokens != nil {
		req["max_tokens"] = *opts.MaxTokens
	}
	if opts.TopP != nil {
		req["top_p"] = *opts.TopP
	}
	if opts.FrequencyPenalty != nil {
		req["frequency_penalty"] = *opts.FrequencyPenalty
	}
	if opts.PresencePenalty != nil {
		req["presence_penalty"] = *opts.PresencePenalty
	}
	if len(opts.Stop) > 0 {
		req["stop"] = opts.Stop
	}
	if len(opts.Tools) > 0 {
		req["tools"] = opts.Tools
	}
	if opts.ToolChoice != nil {
		req["tool_choice"] = opts.ToolChoice
	}
	if opts.ResponseFormat != nil {
		req["response_format"] = opts.ResponseFormat
	}
	if opts.User != "" {
		req["user"] = opts.User
	}
	if opts.Seed != nil {
		req["seed"] = *opts.Seed
	}

	return req
}

// convertMessages converts messages to API format.
func (c *client) convertMessages(messages []Message) []interface{} {
	result := make([]interface{}, len(messages))
	for i, msg := range messages {
		m := map[string]interface{}{
			"role": string(msg.Role),
		}

		// Handle content
		switch content := msg.Content.(type) {
		case TextContent:
			m["content"] = content.Text
		case MultiContent:
			m["content"] = content
		}

		if msg.Name != "" {
			m["name"] = msg.Name
		}
		if len(msg.ToolCalls) > 0 {
			m["tool_calls"] = msg.ToolCalls
		}
		if msg.ToolCallID != "" {
			m["tool_call_id"] = msg.ToolCallID
		}

		result[i] = m
	}
	return result
}

// setHeaders sets authentication and custom headers.
func (c *client) setHeaders(req *http.Request) {
	// Set API key
	if c.apiKey != "" {
		authType := "bearer"
		if c.manifest != nil {
			if t, err := protocol.GetAuthType(c.manifest); err == nil && t != "" {
				authType = t
			}
		}

		switch authType {
		case "bearer":
			req.Header.Set("Authorization", "Bearer "+c.apiKey)
		case "api_key", "header":
			req.Header.Set("Authorization", c.apiKey)
		}
	}

	// Set custom headers
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
}

// parseError parses an error response.
func (c *client) parseError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)
	return &APIError{
		StatusCode: resp.StatusCode,
		Message:    string(body),
	}
}

// APIError represents an API error.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

// Standard error codes (E1001-E9999)
const (
	ErrInvalidRequest  = "E1001"
	ErrAuthentication  = "E1002"
	ErrPermission      = "E1003"
	ErrNotFound        = "E1004"
	ErrRateLimited     = "E2001"
	ErrQuotaExhausted  = "E2002"
	ErrServerOverload  = "E3001"
	ErrInternalError   = "E3002"
	ErrTimeout         = "E3003"
)

// sseStream implements Stream for SSE responses.
type sseStream struct {
	reader   io.ReadCloser
	decoder  *json.Decoder
	event    StreamingEvent
	err      error
	done     bool
	buffer   []byte
}

func newSSEStream(reader io.ReadCloser) *sseStream {
	return &sseStream{
		reader:  reader,
		decoder: json.NewDecoder(reader),
	}
}

func (s *sseStream) Next() bool {
	if s.done {
		return false
	}

	for {
		// Read next line
		line, err := s.readLine()
		if err != nil {
			if err == io.EOF {
				s.done = true
				return false
			}
			s.err = err
			s.done = true
			return false
		}

		// Skip empty lines
		if len(line) == 0 {
			continue
		}

		// Check for data prefix
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")

		// Check for terminal
		if data == "[DONE]" {
			s.done = true
			return false
		}

		// Parse JSON
		var chunk struct {
			Choices []struct {
				Delta struct {
					Content   string     `json:"content"`
					Role      string     `json:"role"`
					ToolCalls []ToolCall `json:"tool_calls"`
				} `json:"delta"`
				FinishReason string `json:"finish_reason"`
			} `json:"choices"`
		}

		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			// Try to continue on parse errors
			continue
		}

		if len(chunk.Choices) > 0 {
			choice := chunk.Choices[0]
			s.event = StreamingEvent{
				Type:         "PartialContentDelta",
				Delta:        choice.Delta.Content,
				FinishReason: choice.FinishReason,
			}
			if len(choice.Delta.ToolCalls) > 0 {
				s.event.ToolCall = &choice.Delta.ToolCalls[0]
			}
			return true
		}
	}
}

func (s *sseStream) readLine() (string, error) {
	for {
		if len(s.buffer) > 0 {
			idx := bytes.IndexByte(s.buffer, '\n')
			if idx >= 0 {
				line := string(s.buffer[:idx])
				s.buffer = s.buffer[idx+1:]
				return strings.TrimSpace(line), nil
			}
		}

		buf := make([]byte, 1024)
		n, err := s.reader.Read(buf)
		if err != nil {
			return "", err
		}
		s.buffer = append(s.buffer, buf[:n]...)
	}
}

func (s *sseStream) Event() StreamingEvent {
	return s.event
}

func (s *sseStream) Err() error {
	return s.err
}

func (s *sseStream) Close() error {
	return s.reader.Close()
}
