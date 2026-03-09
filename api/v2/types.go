// Package v2 provides V2 protocol types for AI-Protocol.
// V2 协议类型定义，采用三环同心圆模型。
package v2

import "time"

// ProviderManifest represents a V2 provider manifest (three-ring model).
// V2 Provider 清单定义，采用三环同心圆模型 (Core/Capabilities/Extensions)。
type ProviderManifest struct {
	Schema          string         `yaml:"$schema" json:"$schema"`
	ID              string         `yaml:"id" json:"id"`
	ProtocolVersion string         `yaml:"protocol_version" json:"protocol_version"`
	Name            string         `yaml:"name" json:"name"`
	Core            CoreRing       `yaml:"core" json:"core"`
	Capabilities    *CapabilitiesRing `yaml:"capabilities,omitempty" json:"capabilities,omitempty"`
	Extensions      *ExtensionsRing   `yaml:"extensions,omitempty" json:"extensions,omitempty"`
	Contracts       []ProviderContract `yaml:"contracts,omitempty" json:"contracts,omitempty"`
}

// CoreRing defines the core ring (L1) - required fields.
// 核心环 (L1) - 必需字段。
type CoreRing struct {
	Endpoint    EndpointConfig    `yaml:"endpoint" json:"endpoint"`
	Auth        AuthConfig        `yaml:"auth" json:"auth"`
	Availability AvailabilityConfig `yaml:"availability" json:"availability"`
	Errors      ErrorConfig       `yaml:"errors" json:"errors"`
}

// EndpointConfig defines endpoint configuration.
// 端点配置。
type EndpointConfig struct {
	BaseURL      string            `yaml:"base_url" json:"base_url"`
	Timeout      time.Duration     `yaml:"timeout" json:"timeout"`
	RetryPolicy  *RetryPolicy      `yaml:"retry_policy" json:"retry_policy"`
	HTTPHeaders  map[string]string `yaml:"http_headers" json:"http_headers"`
}

// AuthConfig defines authentication configuration.
// 认证配置。
type AuthConfig struct {
	Type     string `yaml:"type" json:"type"`         // bearer, api_key, custom
	Location string `yaml:"location" json:"location"` // header, query, body
	Key      string `yaml:"key" json:"key"`           // Authorization, api_key, etc.
	Prefix   string `yaml:"prefix" json:"prefix"`     // Bearer, etc.
}

// AvailabilityConfig defines availability and regions.
// 可用性和区域配置。
type AvailabilityConfig struct {
	Regions    []RegionConfig `yaml:"regions" json:"regions"`
	HealthCheck *HealthCheckConfig `yaml:"health_check" json:"health_check"`
}

// RegionConfig defines a region configuration.
// 区域配置。
type RegionConfig struct {
	ID          string `yaml:"id" json:"id"`
	Name        string `yaml:"name" json:"name"`
	Endpoint    string `yaml:"endpoint" json:"endpoint"`
	Available   bool   `yaml:"available" json:"available"`
}

// HealthCheckConfig defines health check configuration.
// 健康检查配置。
type HealthCheckConfig struct {
	Enabled  bool          `yaml:"enabled" json:"enabled"`
	Interval time.Duration `yaml:"interval" json:"interval"`
	Endpoint string        `yaml:"endpoint" json:"endpoint"`
}

// ErrorConfig defines error handling configuration.
// 错误处理配置。
type ErrorConfig struct {
	StandardCodes []StandardErrorCode `yaml:"standard_codes" json:"standard_codes"`
	Retryable     []string            `yaml:"retryable" json:"retryable"`
}

// StandardErrorCode defines a standard error code mapping.
// 标准错误码映射。
type StandardErrorCode struct {
	Code       string `yaml:"code" json:"code"`             // E1001, E2001, etc.
	HTTPStatus int    `yaml:"http_status" json:"http_status"`
	Message    string `yaml:"message" json:"message"`
}

// CapabilitiesRing defines the capabilities ring (L2).
// 能力环 (L2) - 可选能力声明。
type CapabilitiesRing struct {
	Chat       *ChatCapability       `yaml:"chat,omitempty" json:"chat,omitempty"`
	Streaming  *StreamingCapability  `yaml:"streaming,omitempty" json:"streaming,omitempty"`
	Tools      *ToolsCapability      `yaml:"tools,omitempty" json:"tools,omitempty"`
	Vision     *VisionCapability     `yaml:"vision,omitempty" json:"vision,omitempty"`
	Audio      *AudioCapability      `yaml:"audio,omitempty" json:"audio,omitempty"`
	Video      *VideoCapability      `yaml:"video,omitempty" json:"video,omitempty"`
	Embeddings *EmbeddingsCapability `yaml:"embeddings,omitempty" json:"embeddings,omitempty"`
	Batch      *BatchCapability      `yaml:"batch,omitempty" json:"batch,omitempty"`
	STT        *STTCapability        `yaml:"stt,omitempty" json:"stt,omitempty"`
	TTS        *TTSCapability        `yaml:"tts,omitempty" json:"tts,omitempty"`
	Reranking  *RerankingCapability  `yaml:"reranking,omitempty" json:"reranking,omitempty"`
	Reasoning  *ReasoningCapability  `yaml:"reasoning,omitempty" json:"reasoning,omitempty"`
	MCP        *MCPCapability        `yaml:"mcp,omitempty" json:"mcp,omitempty"`
	ComputerUse *ComputerUseCapability `yaml:"computer_use,omitempty" json:"computer_use,omitempty"`
}

// ChatCapability defines chat capability configuration.
type ChatCapability struct {
	Required       bool     `yaml:"required" json:"required"`
	MaxContext     int      `yaml:"max_context" json:"max_context"`
	SupportedRoles []string `yaml:"supported_roles" json:"supported_roles"`
}

// StreamingCapability defines streaming capability configuration.
type StreamingCapability struct {
	Required   bool                 `yaml:"required" json:"required"`
	Decoder    StreamDecoderConfig  `yaml:"decoder" json:"decoder"`
	EventMap   []EventMappingConfig `yaml:"event_map" json:"event_map"`
	KeepAlive  *KeepAliveConfig     `yaml:"keep_alive" json:"keep_alive"`
}

// StreamDecoderConfig defines stream decoder configuration.
type StreamDecoderConfig struct {
	Format   string `yaml:"format" json:"format"`
	Strategy string `yaml:"strategy" json:"strategy"`
}

// EventMappingConfig defines event mapping configuration.
type EventMappingConfig struct {
	Match   string            `yaml:"match" json:"match"`
	Emit    string            `yaml:"emit" json:"emit"`
	Extract map[string]string `yaml:"extract" json:"extract"`
}

// KeepAliveConfig defines keep-alive configuration.
type KeepAliveConfig struct {
	Enabled     bool          `yaml:"enabled" json:"enabled"`
	Interval    time.Duration `yaml:"interval" json:"interval"`
	IdleTimeout time.Duration `yaml:"idle_timeout" json:"idle_timeout"`
}

// ToolsCapability defines tools capability configuration.
type ToolsCapability struct {
	Required          bool     `yaml:"required" json:"required"`
	ParallelToolCalls bool     `yaml:"parallel_tool_calls" json:"parallel_tool_calls"`
	MaxTools          int      `yaml:"max_tools" json:"max_tools"`
	SupportedTypes    []string `yaml:"supported_types" json:"supported_types"`
}

// VisionCapability defines vision capability configuration.
type VisionCapability struct {
	Required        bool     `yaml:"required" json:"required"`
	SupportedFormats []string `yaml:"supported_formats" json:"supported_formats"`
	MaxSize         int      `yaml:"max_size" json:"max_size"`
}

// AudioCapability defines audio capability configuration.
type AudioCapability struct {
	Required        bool     `yaml:"required" json:"required"`
	SupportedFormats []string `yaml:"supported_formats" json:"supported_formats"`
	MaxDuration     int      `yaml:"max_duration" json:"max_duration"`
}

// VideoCapability defines video capability configuration.
type VideoCapability struct {
	Required        bool     `yaml:"required" json:"required"`
	SupportedFormats []string `yaml:"supported_formats" json:"supported_formats"`
	MaxDuration     int      `yaml:"max_duration" json:"max_duration"`
	Generation      *VideoGenerationConfig `yaml:"generation,omitempty" json:"generation,omitempty"`
}

// VideoGenerationConfig defines video generation configuration.
type VideoGenerationConfig struct {
	AsyncPolling    bool     `yaml:"async_polling" json:"async_polling"`
	PollInterval    time.Duration `yaml:"poll_interval" json:"poll_interval"`
	MaxPollAttempts int      `yaml:"max_poll_attempts" json:"max_poll_attempts"`
}

// EmbeddingsCapability defines embeddings capability configuration.
type EmbeddingsCapability struct {
	Required      bool     `yaml:"required" json:"required"`
	Dimensions    []int    `yaml:"dimensions" json:"dimensions"`
	BatchSize     int      `yaml:"batch_size" json:"batch_size"`
}

// BatchCapability defines batch capability configuration.
type BatchCapability struct {
	Required       bool     `yaml:"required" json:"required"`
	MaxBatchSize   int      `yaml:"max_batch_size" json:"max_batch_size"`
	RetentionDays  int      `yaml:"retention_days" json:"retention_days"`
}

// STTCapability defines speech-to-text capability configuration.
type STTCapability struct {
	Required        bool     `yaml:"required" json:"required"`
	SupportedFormats []string `yaml:"supported_formats" json:"supported_formats"`
	Languages       []string `yaml:"languages" json:"languages"`
}

// TTSCapability defines text-to-speech capability configuration.
type TTSCapability struct {
	Required       bool     `yaml:"required" json:"required"`
	Voices         []string `yaml:"voices" json:"voices"`
	SpeedRange     []float64 `yaml:"speed_range" json:"speed_range"`
}

// RerankingCapability defines reranking capability configuration.
type RerankingCapability struct {
	Required       bool     `yaml:"required" json:"required"`
	MaxDocuments   int      `yaml:"max_documents" json:"max_documents"`
	ReturnDocuments bool     `yaml:"return_documents" json:"return_documents"`
}

// ReasoningCapability defines reasoning capability configuration.
type ReasoningCapability struct {
	Required       bool     `yaml:"required" json:"required"`
	MaxThinkingTokens int   `yaml:"max_thinking_tokens" json:"max_thinking_tokens"`
}

// MCPCapability defines MCP capability configuration.
type MCPCapability struct {
	Required      bool     `yaml:"required" json:"required"`
	Transport     string   `yaml:"transport" json:"transport"`
	ProtocolLevel string   `yaml:"protocol_level" json:"protocol_level"`
}

// ComputerUseCapability defines computer use capability configuration.
type ComputerUseCapability struct {
	Required     bool     `yaml:"required" json:"required"`
	MaxActions   int      `yaml:"max_actions" json:"max_actions"`
	SupportedOS  []string `yaml:"supported_os" json:"supported_os"`
}

// ExtensionsRing defines the extensions ring (L3).
// 扩展环 (L3) - 环境和高级功能。
type ExtensionsRing struct {
	Multimodal     *MultimodalExtension     `yaml:"multimodal,omitempty" json:"multimodal,omitempty"`
	ContextPolicy  *ContextPolicyExtension  `yaml:"context_policy,omitempty" json:"context_policy,omitempty"`
	RateLimit      *RateLimitExtension      `yaml:"rate_limit,omitempty" json:"rate_limit,omitempty"`
	Custom         map[string]interface{}   `yaml:"custom,omitempty" json:"custom,omitempty"`
}

// MultimodalExtension defines multimodal extension configuration.
type MultimodalExtension struct {
	Interleaving string `yaml:"interleaving" json:"interleaving"`
	MaxInputs    int    `yaml:"max_inputs" json:"max_inputs"`
}

// ContextPolicyExtension defines context policy configuration.
type ContextPolicyExtension struct {
	MaxTokens      int    `yaml:"max_tokens" json:"max_tokens"`
	TruncationStrategy string `yaml:"truncation_strategy" json:"truncation_strategy"`
}

// RateLimitExtension defines rate limit configuration.
type RateLimitExtension struct {
	RequestsPerMinute int `yaml:"requests_per_minute" json:"requests_per_minute"`
	TokensPerMinute   int `yaml:"tokens_per_minute" json:"tokens_per_minute"`
}

// ProviderContract defines a provider contract.
// Provider 契约定义。
type ProviderContract struct {
	Version      string            `yaml:"version" json:"version"`
	APIFamily    string            `yaml:"api_family" json:"api_family"`
	RequestShape map[string]string `yaml:"request_shape" json:"request_shape"`
	ResponseShape map[string]string `yaml:"response_shape" json:"response_shape"`
}

// RetryPolicy defines retry behavior.
type RetryPolicy struct {
	Strategy      string `yaml:"strategy" json:"strategy"`
	MinDelay      int    `yaml:"min_delay_ms" json:"min_delay_ms"`
	MaxDelay      int    `yaml:"max_delay_ms" json:"max_delay_ms"`
	MaxAttempts   int    `yaml:"max_attempts" json:"max_attempts"`
	Jitter        string `yaml:"jitter" json:"jitter"`
	RetryOnStatus []int  `yaml:"retry_on_status" json:"retry_on_status"`
}
