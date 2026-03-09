package ailib

import (
	"context"
	"encoding/json"
	"fmt"
)

// EmbeddingRequest represents an embedding request.
type EmbeddingRequest struct {
	Model          string   `json:"model"`
	Input          []string `json:"input"`
	EncodingFormat string   `json:"encoding_format,omitempty"` // float, base64
	Dimensions     int      `json:"dimensions,omitempty"`
	User           string   `json:"user,omitempty"`
}

// EmbeddingResponse represents an embedding response.
type EmbeddingResponse struct {
	Object string             `json:"object"`
	Data   []EmbeddingData    `json:"data"`
	Model  string             `json:"model"`
	Usage  *EmbeddingUsage    `json:"usage,omitempty"`
}

// EmbeddingData represents a single embedding.
type EmbeddingData struct {
	Object    string    `json:"object"`
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
}

// EmbeddingUsage represents embedding token usage.
type EmbeddingUsage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// EmbeddingsClient provides embeddings API.
type EmbeddingsClient struct {
	client *client
}

// NewEmbeddingsClient creates an embeddings client.
func NewEmbeddingsClient(c Client) (*EmbeddingsClient, error) {
	ac, ok := c.(*client)
	if !ok {
		return nil, fmt.Errorf("invalid client type")
	}
	return &EmbeddingsClient{client: ac}, nil
}

// Create creates embeddings.
func (e *EmbeddingsClient) Create(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error) {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := e.client.baseURL + "/embeddings"
	httpReq, err := newRequest(ctx, "POST", url, reqBytes)
	if err != nil {
		return nil, err
	}

	e.client.setHeaders(httpReq)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := e.client.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, e.client.parseError(resp)
	}

	var response EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

func newRequest(ctx context.Context, method, url string, body []byte) (*http.Request, error) {
	return newRequestWithBody(ctx, method, url, body)
}

func newRequestWithBody(ctx context.Context, method, url string, body []byte) (*http.Request, error) {
	var bodyReader interface{}
	if len(body) > 0 {
		bodyReader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	return req, nil
}
