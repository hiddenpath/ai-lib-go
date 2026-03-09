package ailib

import (
	"context"
	"encoding/json"
	"fmt"
)

// RerankRequest represents a reranking request.
type RerankRequest struct {
	Model           string          `json:"model"`
	Query           string          `json:"query"`
	Documents       []RerankDocument `json:"documents"`
	TopN            int             `json:"top_n,omitempty"`
	ReturnDocuments bool            `json:"return_documents,omitempty"`
	MaxChunksPerDoc int             `json:"max_chunks_per_doc,omitempty"`
}

// RerankDocument represents a document for reranking.
type RerankDocument struct {
	Text string `json:"text"`
}

// RerankResponse represents a reranking response.
type RerankResponse struct {
	Object string       `json:"object"`
	Data   []RerankResult `json:"data"`
	Model  string       `json:"model"`
	Usage  *RerankUsage `json:"usage,omitempty"`
}

// RerankResult represents a single rerank result.
type RerankResult struct {
	Index          int             `json:"index"`
	RelevanceScore float64         `json:"relevance_score"`
	Document       *RerankDocument `json:"document,omitempty"`
}

// RerankUsage represents reranking token usage.
type RerankUsage struct {
	TotalTokens int `json:"total_tokens"`
}

// RerankClient provides reranking API.
type RerankClient struct {
	client *client
}

// NewRerankClient creates a rerank client.
func NewRerankClient(c Client) (*RerankClient, error) {
	ac, ok := c.(*client)
	if !ok {
		return nil, fmt.Errorf("invalid client type")
	}
	return &RerankClient{client: ac}, nil
}

// Rerank reranks documents by relevance to query.
func (r *RerankClient) Rerank(ctx context.Context, req *RerankRequest) (*RerankResponse, error) {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := r.client.baseURL + "/rerank"
	httpReq, err := newRequest(ctx, "POST", url, reqBytes)
	if err != nil {
		return nil, err
	}

	r.client.setHeaders(httpReq)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := r.client.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, r.client.parseError(resp)
	}

	var response RerankResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}
