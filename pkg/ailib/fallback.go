// Package ailib fallback execution helpers.
// 回退执行器，提供最小可用的跨客户端失败转移能力。
package ailib

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type FallbackClient struct {
	clients []Client
	policy  FallbackPolicy
	mu      sync.Mutex
	states  []fallbackState
}

type FallbackPolicy struct {
	FailureThreshold int
	CircuitOpenFor   time.Duration
}

type fallbackState struct {
	consecutiveFailures int
	circuitOpenUntil    time.Time
}

type ProviderHealth struct {
	Index               int
	ConsecutiveFailures int
	CircuitOpen         bool
	CircuitOpenUntil    time.Time
}

func DefaultFallbackPolicy() FallbackPolicy {
	return FallbackPolicy{
		FailureThreshold: 3,
		CircuitOpenFor:   30 * time.Second,
	}
}

func NewFallbackClient(primary Client, fallbacks ...Client) *FallbackClient {
	return NewFallbackClientWithPolicy(DefaultFallbackPolicy(), primary, fallbacks...)
}

func NewFallbackClientWithPolicy(policy FallbackPolicy, primary Client, fallbacks ...Client) *FallbackClient {
	if policy.FailureThreshold <= 0 {
		policy.FailureThreshold = 1
	}
	if policy.CircuitOpenFor <= 0 {
		policy.CircuitOpenFor = 5 * time.Second
	}
	clients := make([]Client, 0, 1+len(fallbacks))
	if primary != nil {
		clients = append(clients, primary)
	}
	for _, c := range fallbacks {
		if c != nil {
			clients = append(clients, c)
		}
	}
	return &FallbackClient{
		clients: clients,
		policy:  policy,
		states:  make([]fallbackState, len(clients)),
	}
}

func (f *FallbackClient) Chat(ctx context.Context, messages []Message, opts *ChatOptions) (*ChatResponse, error) {
	return callWithFallback(f, func(c Client) (*ChatResponse, error) {
		return c.Chat(ctx, messages, opts)
	})
}

func (f *FallbackClient) ChatStream(ctx context.Context, messages []Message, opts *ChatOptions) (Stream, error) {
	return callWithFallback(f, func(c Client) (Stream, error) {
		return c.ChatStream(ctx, messages, opts)
	})
}

func (f *FallbackClient) Embeddings(ctx context.Context, req EmbeddingRequest) (*EmbeddingResponse, error) {
	return callWithFallback(f, func(c Client) (*EmbeddingResponse, error) {
		return c.Embeddings(ctx, req)
	})
}

func (f *FallbackClient) BatchCreate(ctx context.Context, req BatchCreateRequest) (*BatchJob, error) {
	return callWithFallback(f, func(c Client) (*BatchJob, error) {
		return c.BatchCreate(ctx, req)
	})
}

func (f *FallbackClient) BatchGet(ctx context.Context, batchID string) (*BatchJob, error) {
	return callWithFallback(f, func(c Client) (*BatchJob, error) {
		return c.BatchGet(ctx, batchID)
	})
}

func (f *FallbackClient) BatchCancel(ctx context.Context, batchID string) (*BatchJob, error) {
	return callWithFallback(f, func(c Client) (*BatchJob, error) {
		return c.BatchCancel(ctx, batchID)
	})
}

func (f *FallbackClient) STTTranscribe(ctx context.Context, req STTRequest) (*STTResponse, error) {
	return callWithFallback(f, func(c Client) (*STTResponse, error) {
		return c.STTTranscribe(ctx, req)
	})
}

func (f *FallbackClient) TTSSpeak(ctx context.Context, req TTSRequest) (*TTSResponse, error) {
	return callWithFallback(f, func(c Client) (*TTSResponse, error) {
		return c.TTSSpeak(ctx, req)
	})
}

func (f *FallbackClient) Rerank(ctx context.Context, req RerankRequest) (*RerankResponse, error) {
	return callWithFallback(f, func(c Client) (*RerankResponse, error) {
		return c.Rerank(ctx, req)
	})
}

func (f *FallbackClient) MCPListTools(ctx context.Context) (*MCPListToolsResponse, error) {
	return callWithFallback(f, func(c Client) (*MCPListToolsResponse, error) {
		return c.MCPListTools(ctx)
	})
}

func (f *FallbackClient) MCPCallTool(ctx context.Context, req MCPCallToolRequest) (*MCPCallToolResponse, error) {
	return callWithFallback(f, func(c Client) (*MCPCallToolResponse, error) {
		return c.MCPCallTool(ctx, req)
	})
}

func (f *FallbackClient) ComputerUse(ctx context.Context, req ComputerUseRequest) (*ComputerUseResponse, error) {
	return callWithFallback(f, func(c Client) (*ComputerUseResponse, error) {
		return c.ComputerUse(ctx, req)
	})
}

func (f *FallbackClient) Reason(ctx context.Context, req ReasoningRequest) (*ReasoningResponse, error) {
	return callWithFallback(f, func(c Client) (*ReasoningResponse, error) {
		return c.Reason(ctx, req)
	})
}

func (f *FallbackClient) VideoGenerate(ctx context.Context, req VideoGenerateRequest) (*VideoJob, error) {
	return callWithFallback(f, func(c Client) (*VideoJob, error) {
		return c.VideoGenerate(ctx, req)
	})
}

func (f *FallbackClient) VideoGet(ctx context.Context, jobID string) (*VideoJob, error) {
	return callWithFallback(f, func(c Client) (*VideoJob, error) {
		return c.VideoGet(ctx, jobID)
	})
}

func (f *FallbackClient) Close() error {
	var firstErr error
	for _, c := range f.clients {
		if err := c.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (f *FallbackClient) HealthSnapshot() []ProviderHealth {
	now := time.Now()
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([]ProviderHealth, 0, len(f.states))
	for i, s := range f.states {
		out = append(out, ProviderHealth{
			Index:               i,
			ConsecutiveFailures: s.consecutiveFailures,
			CircuitOpen:         s.circuitOpenUntil.After(now),
			CircuitOpenUntil:    s.circuitOpenUntil,
		})
	}
	return out
}

func callWithFallback[T any](f *FallbackClient, fn func(Client) (T, error)) (T, error) {
	var zero T
	if len(f.clients) == 0 {
		return zero, fmt.Errorf("no clients configured for fallback execution")
	}
	var lastErr error
	attempted := 0
	for i, c := range f.clients {
		if !f.canAttempt(i, time.Now()) {
			continue
		}
		attempted++
		resp, err := fn(c)
		if err == nil {
			f.recordSuccess(i)
			return resp, nil
		}
		lastErr = err
		f.recordFailure(i)
		if i == len(f.clients)-1 || !isFallbackableErr(err) {
			return zero, err
		}
	}
	if attempted == 0 {
		return zero, fmt.Errorf("all providers are temporarily unhealthy (circuits open)")
	}
	return zero, lastErr
}

func (f *FallbackClient) canAttempt(idx int, now time.Time) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	if idx >= len(f.states) {
		return true
	}
	return !f.states[idx].circuitOpenUntil.After(now)
}

func (f *FallbackClient) recordSuccess(idx int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if idx >= len(f.states) {
		return
	}
	f.states[idx].consecutiveFailures = 0
	f.states[idx].circuitOpenUntil = time.Time{}
}

func (f *FallbackClient) recordFailure(idx int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	if idx >= len(f.states) {
		return
	}
	s := &f.states[idx]
	s.consecutiveFailures++
	if s.consecutiveFailures >= f.policy.FailureThreshold {
		s.circuitOpenUntil = time.Now().Add(f.policy.CircuitOpenFor)
	}
}
