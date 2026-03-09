// Package v1 provides V1 protocol types for AI-Protocol.
// V1 协议类型定义，用于兼容现有的 provider manifests。
package v1

import "time"

// ProviderManifest represents a V1 provider manifest.
// Provider 清单定义，包含 provider 的完整配置。
type ProviderManifest struct {
	Schema         string          `yaml:"$schema" json:"$schema"`
	ID             string          `yaml:"id" json:"id"`
	ProtocolVersion string         `yaml:"protocol_version" json:"protocol_version"`
	Name           string          `yaml:"name" json:"name"`
	BaseURL        string          `yaml:"base_url" json:"base_url"`
	APIFormat      string          `yaml:"api_format" json:"api_format"` // openai, anthropic, gemini
	Auth           AuthConfig      `yaml:"auth" json:"auth"`
	Endpoints      []Endpoint      `yaml:"endpoints" json:"endpoints"`
	Streaming      *StreamingConfig `yaml:"streaming" json:"streaming"`
	ErrorClassification ErrorClassification `yaml:"error_classification" json:"error_classification"`
	RateLimitHeaders  *RateLimitHeaders `yaml:"rate_limit_headers" json:"rate_limit_headers"`
	RetryPolicy       *RetryPolicy     `yaml:"retry_policy" json:"retry_policy"`
	TerminationReasons []TerminationReason `yaml:"termination_reasons" json:"termination_reasons"`
	APIFamilies      []string         `yaml:"api_families" json:"api_families"`
}

// AuthConfig defines authentication configuration.
// 认证配置，定义 API key 的传递方式。
type AuthConfig struct {
	Type   string `yaml:"type" json:"type"`     // bearer, header, query
	Header string `yaml:"header" json:"header"` // Authorization, X-API-Key, etc.
	Prefix string `yaml:"prefix" json:"prefix"` // Bearer, Token, etc.
	Query  string `yaml:"query" json:"query"`   // api_key, key, etc.
}

// Endpoint defines an API endpoint.
// API 端点定义，包含路径和 HTTP 方法。
type Endpoint struct {
	Name           string            `yaml:"name" json:"name"`
	Path           string            `yaml:"path" json:"path"`
	Method         string            `yaml:"method" json:"method"`
	OperationID    string            `yaml:"operation_id" json:"operation_id"`
	Params         map[string]string `yaml:"params" json:"params"`
	RequestFormat  string            `yaml:"request_format" json:"request_format"`
	ResponseFormat string            `yaml:"response_format" json:"response_format"`
}

// StreamingConfig defines streaming configuration.
// 流式响应配置，定义 SSE 解码和事件映射。
type StreamingConfig struct {
	Decoder   StreamDecoder `yaml:"decoder" json:"decoder"`
	EventMap  []EventMapping `yaml:"event_map" json:"event_map"`
	KeepAlive *KeepAliveConfig `yaml:"keep_alive" json:"keep_alive"`
}

// StreamDecoder defines how to decode streaming responses.
// 流解码器配置。
type StreamDecoder struct {
	Format   string `yaml:"format" json:"format"`     // openai_sse, anthropic_sse, gemini_sse
	Strategy string `yaml:"strategy" json:"strategy"` // line_based, event_stream, json_lines
	Delim    string `yaml:"delim" json:"delim"`       // data: , event: , etc.
	Terminal string `yaml:"terminal" json:"terminal"` // [DONE], etc.
}

// EventMapping defines how to map provider events to standard events.
// 事件映射规则，将 provider 特定事件转换为标准事件。
type EventMapping struct {
	Match   string            `yaml:"match" json:"match"`     // JSONPath expression
	Emit    string            `yaml:"emit" json:"emit"`       // Event type to emit
	Extract map[string]string `yaml:"extract" json:"extract"` // Field extraction rules
}

// KeepAliveConfig defines keep-alive behavior for streaming.
// Keep-alive 配置，用于保持长连接。
type KeepAliveConfig struct {
	Enabled     bool          `yaml:"enabled" json:"enabled"`
	Interval    time.Duration `yaml:"interval" json:"interval"`
	IdleTimeout time.Duration `yaml:"idle_timeout" json:"idle_timeout"`
}

// ErrorClassification maps HTTP status codes to error types.
// 错误分类，将 HTTP 状态码映射到标准错误类型。
type ErrorClassification struct {
	ByHTTPStatus map[string]string `yaml:"by_http_status" json:"by_http_status"`
	ByErrorCode  map[string]string `yaml:"by_error_code" json:"by_error_code"`
}

// RateLimitHeaders defines rate limit header names.
// Rate limit 头定义。
type RateLimitHeaders struct {
	RequestsLimit    string `yaml:"requests_limit" json:"requests_limit"`
	RequestsRemaining string `yaml:"requests_remaining" json:"requests_remaining"`
	RetryAfter       string `yaml:"retry_after" json:"retry_after"`
}

// RetryPolicy defines retry behavior.
// 重试策略配置。
type RetryPolicy struct {
	Strategy       string `yaml:"strategy" json:"strategy"`           // exponential_backoff, fixed
	MinDelay       int    `yaml:"min_delay_ms" json:"min_delay_ms"`   // Minimum delay in ms
	MaxDelay       int    `yaml:"max_delay_ms" json:"max_delay_ms"`   // Maximum delay in ms
	MaxAttempts    int    `yaml:"max_attempts" json:"max_attempts"`   // Max retry attempts
	Jitter         string `yaml:"jitter" json:"jitter"`               // none, full, decorrelated
	RetryOnStatus  []int  `yaml:"retry_on_http_status" json:"retry_on_http_status"`
	RetryOnErrors  []string `yaml:"retry_on_errors" json:"retry_on_errors"`
}

// TerminationReason defines termination reason mapping.
// 终止原因映射。
type TerminationReason struct {
	ProviderReason string `yaml:"provider_reason" json:"provider_reason"`
	StandardReason string `yaml:"standard_reason" json:"standard_reason"`
}

// ModelManifest represents a V1 model manifest.
// Model 清单定义。
type ModelManifest struct {
	Schema string `yaml:"$schema" json:"$schema"`
	Models map[string]Model `yaml:"models" json:"models"`
}

// Model represents a model configuration.
// 模型配置。
type Model struct {
	Provider      string   `yaml:"provider" json:"provider"`
	ModelID       string   `yaml:"model_id" json:"model_id"`
	ContextWindow int      `yaml:"context_window" json:"context_window"`
	Capabilities  []string `yaml:"capabilities" json:"capabilities"`
	Pricing       *Pricing `yaml:"pricing" json:"pricing"`
	Aliases       []string `yaml:"aliases" json:"aliases"`
}

// Pricing represents model pricing information.
// 模型定价信息。
type Pricing struct {
	InputPerToken  float64 `yaml:"input_per_token" json:"input_per_token"`
	OutputPerToken float64 `yaml:"output_per_token" json:"output_per_token"`
}
