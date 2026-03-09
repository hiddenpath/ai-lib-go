package ailib

import (
	"context"
	"encoding/json"
	"fmt"
)

// BatchRequest represents a batch processing request.
type BatchRequest struct {
	CustomID string      `json:"custom_id"`
	Method   string      `json:"method"`
	URL      string      `json:"url"`
	Body     interface{} `json:"body"`
}

// BatchResponse represents a batch processing response.
type BatchResponse struct {
	CustomID string          `json:"custom_id"`
	Response *BatchResult    `json:"response"`
	Error    *BatchError     `json:"error,omitempty"`
}

// BatchResult represents a successful batch result.
type BatchResult struct {
	StatusCode int             `json:"status_code"`
	Body       json.RawMessage `json:"body"`
}

// BatchError represents a batch error.
type BatchError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// BatchJob represents a batch job.
type BatchJob struct {
	ID              string        `json:"id"`
	Object          string        `json:"object"`
	Status          string        `json:"status"` // pending, running, completed, failed
	InputFileID     string        `json:"input_file_id"`
	OutputFileID    string        `json:"output_file_id,omitempty"`
	ErrorFileID     string        `json:"error_file_id,omitempty"`
	CreatedAt       int64         `json:"created_at"`
	CompletedAt     int64         `json:"completed_at,omitempty"`
	ExpiresAt       int64         `json:"expires_at,omitempty"`
	RequestCounts   *RequestCounts `json:"request_counts,omitempty"`
	Metadata        map[string]any `json:"metadata,omitempty"`
}

// RequestCounts represents batch request counts.
type RequestCounts struct {
	Total     int `json:"total"`
	Completed int `json:"completed"`
	Failed    int `json:"failed"`
}

// BatchClient provides batch processing API.
type BatchClient struct {
	client *client
}

// NewBatchClient creates a batch client.
func NewBatchClient(c Client) (*BatchClient, error) {
	ac, ok := c.(*client)
	if !ok {
		return nil, fmt.Errorf("invalid client type")
	}
	return &BatchClient{client: ac}, nil
}

// Create creates a new batch job.
func (b *BatchClient) Create(ctx context.Context, inputFileID string, endpoint string, completionWindow string) (*BatchJob, error) {
	req := map[string]interface{}{
		"input_file_id":       inputFileID,
		"endpoint":            endpoint,
		"completion_window":   completionWindow,
	}
	reqBytes, _ := json.Marshal(req)

	url := b.client.baseURL + "/batches"
	httpReq, err := newRequest(ctx, "POST", url, reqBytes)
	if err != nil {
		return nil, err
	}

	b.client.setHeaders(httpReq)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := b.client.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, b.client.parseError(resp)
	}

	var batch BatchJob
	if err := json.NewDecoder(resp.Body).Decode(&batch); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &batch, nil
}

// Retrieve retrieves a batch job.
func (b *BatchClient) Retrieve(ctx context.Context, batchID string) (*BatchJob, error) {
	url := b.client.baseURL + "/batches/" + batchID
	httpReq, err := newRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	b.client.setHeaders(httpReq)

	resp, err := b.client.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, b.client.parseError(resp)
	}

	var batch BatchJob
	if err := json.NewDecoder(resp.Body).Decode(&batch); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &batch, nil
}

// Cancel cancels a batch job.
func (b *BatchClient) Cancel(ctx context.Context, batchID string) (*BatchJob, error) {
	url := b.client.baseURL + "/batches/" + batchID + "/cancel"
	httpReq, err := newRequest(ctx, "POST", url, nil)
	if err != nil {
		return nil, err
	}

	b.client.setHeaders(httpReq)

	resp, err := b.client.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, b.client.parseError(resp)
	}

	var batch BatchJob
	if err := json.NewDecoder(resp.Body).Decode(&batch); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &batch, nil
}
