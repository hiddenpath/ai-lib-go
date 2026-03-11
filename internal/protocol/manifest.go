// Package protocol defines protocol manifest contracts.
// 协议清单模型，统一 V1/V2 读取与端点解析。
package protocol

import "fmt"

type V1Manifest struct {
	ID              string           `yaml:"id" json:"id"`
	ProtocolVersion string           `yaml:"protocol_version" json:"protocol_version"`
	BaseURL         string           `yaml:"base_url" json:"base_url"`
	APIFormat       string           `yaml:"api_format" json:"api_format"`
	Capabilities    []string         `yaml:"capabilities" json:"capabilities"`
	ErrorClass      ErrorClass       `yaml:"error_classification" json:"error_classification"`
	RetryPolicy     RetryPolicy      `yaml:"retry_policy" json:"retry_policy"`
	Auth            V1Auth           `yaml:"auth" json:"auth"`
	Endpoints       []V1Endpoint     `yaml:"endpoints" json:"endpoints"`
	Streaming       *StreamingConfig `yaml:"streaming" json:"streaming"`
}

type StreamingConfig struct {
	Decoder *DecoderConfig `yaml:"decoder" json:"decoder"`
}

type DecoderConfig struct {
	Format   string `yaml:"format" json:"format"`
	Strategy string `yaml:"strategy" json:"strategy"`
}

type V1Auth struct {
	Type   string `yaml:"type" json:"type"`
	Header string `yaml:"header" json:"header"`
	Prefix string `yaml:"prefix" json:"prefix"`
}

type V1Endpoint struct {
	Name        string `yaml:"name" json:"name"`
	OperationID string `yaml:"operation_id" json:"operation_id"`
	Path        string `yaml:"path" json:"path"`
	Method      string `yaml:"method" json:"method"`
}

type V2Manifest struct {
	ID              string           `yaml:"id" json:"id"`
	ProtocolVersion string           `yaml:"protocol_version" json:"protocol_version"`
	Core            V2Core           `yaml:"core" json:"core"`
	ErrorClass      ErrorClass       `yaml:"error_classification" json:"error_classification"`
	RetryPolicy     RetryPolicy      `yaml:"retry_policy" json:"retry_policy"`
	Capabilities    V2Caps           `yaml:"capabilities" json:"capabilities"`
	Streaming       *StreamingConfig `yaml:"streaming" json:"streaming"`
}

type V2Core struct {
	Endpoint V2Endpoint `yaml:"endpoint" json:"endpoint"`
	Auth     V2Auth     `yaml:"auth" json:"auth"`
}

type V2Endpoint struct {
	BaseURL   string                 `yaml:"base_url" json:"base_url"`
	Endpoints map[string]V2Operation `yaml:"endpoints" json:"endpoints"`
}

type V2Operation struct {
	Path   string `yaml:"path" json:"path"`
	Method string `yaml:"method" json:"method"`
}

type V2Auth struct {
	Type   string `yaml:"type" json:"type"`
	Key    string `yaml:"key" json:"key"`
	Prefix string `yaml:"prefix" json:"prefix"`
}

type V2Caps struct {
	Chat        *RequiredCap `yaml:"chat" json:"chat"`
	Streaming   *RequiredCap `yaml:"streaming" json:"streaming"`
	Tools       *RequiredCap `yaml:"tools" json:"tools"`
	Vision      *RequiredCap `yaml:"vision" json:"vision"`
	Audio       *RequiredCap `yaml:"audio" json:"audio"`
	Video       *RequiredCap `yaml:"video" json:"video"`
	Embeddings  *RequiredCap `yaml:"embeddings" json:"embeddings"`
	Batch       *RequiredCap `yaml:"batch" json:"batch"`
	STT         *RequiredCap `yaml:"stt" json:"stt"`
	TTS         *RequiredCap `yaml:"tts" json:"tts"`
	Reranking   *RequiredCap `yaml:"reranking" json:"reranking"`
	MCP         *RequiredCap `yaml:"mcp" json:"mcp"`
	ComputerUse *RequiredCap `yaml:"computer_use" json:"computer_use"`
	Reasoning   *RequiredCap `yaml:"reasoning" json:"reasoning"`
}

type RequiredCap struct {
	Required bool `yaml:"required" json:"required"`
}

type ErrorClass struct {
	ByHTTPStatus map[string]string `yaml:"by_http_status" json:"by_http_status"`
	ByErrorCode  map[string]string `yaml:"by_error_code" json:"by_error_code"`
	ByErrorType  map[string]string `yaml:"by_error_type" json:"by_error_type"`
}

type RetryPolicy struct {
	MaxRetries int `yaml:"max_retries" json:"max_retries"`
}

func BaseURL(m any) (string, error) {
	switch v := m.(type) {
	case *V1Manifest:
		return v.BaseURL, nil
	case *V2Manifest:
		return v.Core.Endpoint.BaseURL, nil
	default:
		return "", fmt.Errorf("unsupported manifest type: %T", m)
	}
}

func AuthHeader(m any) (name string, valuePrefix string, err error) {
	switch v := m.(type) {
	case *V1Manifest:
		h := v.Auth.Header
		if h == "" {
			h = "Authorization"
		}
		p := v.Auth.Prefix
		if p == "" && v.Auth.Type == "bearer" {
			p = "Bearer "
		}
		return h, p, nil
	case *V2Manifest:
		h := v.Core.Auth.Key
		if h == "" {
			h = "Authorization"
		}
		p := v.Core.Auth.Prefix
		if p == "" && v.Core.Auth.Type == "bearer" {
			p = "Bearer "
		}
		return h, p, nil
	default:
		return "", "", fmt.Errorf("unsupported manifest type: %T", m)
	}
}

func EndpointFor(m any, key string, fallback string) (path string, method string) {
	switch v := m.(type) {
	case *V1Manifest:
		for _, ep := range v.Endpoints {
			if ep.OperationID == key || ep.Name == key {
				return ep.Path, upperOrPOST(ep.Method)
			}
		}
	case *V2Manifest:
		if op, ok := v.Core.Endpoint.Endpoints[key]; ok {
			return op.Path, upperOrPOST(op.Method)
		}
	}
	return fallback, "POST"
}

func HasCapability(m any, name string) bool {
	switch v := m.(type) {
	case *V1Manifest:
		for _, c := range v.Capabilities {
			if c == name {
				return true
			}
		}
		return false
	case *V2Manifest:
		switch name {
		case "chat":
			return v.Capabilities.Chat != nil
		case "streaming":
			return v.Capabilities.Streaming != nil
		case "tools":
			return v.Capabilities.Tools != nil
		case "vision":
			return v.Capabilities.Vision != nil
		case "audio":
			return v.Capabilities.Audio != nil
		case "video":
			return v.Capabilities.Video != nil
		case "embeddings":
			return v.Capabilities.Embeddings != nil
		case "batch":
			return v.Capabilities.Batch != nil
		case "stt":
			return v.Capabilities.STT != nil
		case "tts":
			return v.Capabilities.TTS != nil
		case "reranking":
			return v.Capabilities.Reranking != nil
		case "mcp":
			return v.Capabilities.MCP != nil
		case "computer_use":
			return v.Capabilities.ComputerUse != nil
		case "reasoning":
			return v.Capabilities.Reasoning != nil
		default:
			return false
		}
	default:
		return true
	}
}

func ClassifyError(m any, status int, providerCode, providerType string) (code string, ok bool) {
	cls, exists := errorClass(m)
	if !exists {
		return "", false
	}
	if providerCode != "" && cls.ByErrorCode != nil {
		if n, hit := cls.ByErrorCode[providerCode]; hit {
			return normalizeErrorNameToCode(n)
		}
	}
	if providerType != "" && cls.ByErrorType != nil {
		if n, hit := cls.ByErrorType[providerType]; hit {
			return normalizeErrorNameToCode(n)
		}
	}
	if cls.ByHTTPStatus != nil {
		if n, hit := cls.ByHTTPStatus[fmt.Sprintf("%d", status)]; hit {
			return normalizeErrorNameToCode(n)
		}
	}
	return "", false
}

func RetryMaxAttempts(m any) (int, bool) {
	switch v := m.(type) {
	case *V1Manifest:
		if v.RetryPolicy.MaxRetries > 0 {
			return v.RetryPolicy.MaxRetries + 1, true
		}
	case *V2Manifest:
		if v.RetryPolicy.MaxRetries > 0 {
			return v.RetryPolicy.MaxRetries + 1, true
		}
	}
	return 0, false
}

func errorClass(m any) (ErrorClass, bool) {
	switch v := m.(type) {
	case *V1Manifest:
		return v.ErrorClass, true
	case *V2Manifest:
		return v.ErrorClass, true
	default:
		return ErrorClass{}, false
	}
}

func normalizeErrorNameToCode(name string) (string, bool) {
	switch name {
	case "invalid_request":
		return "E1001", true
	case "authentication":
		return "E1002", true
	case "permission_denied":
		return "E1003", true
	case "not_found":
		return "E1004", true
	case "request_too_large":
		return "E1005", true
	case "rate_limited":
		return "E2001", true
	case "quota_exhausted":
		return "E2002", true
	case "server_error":
		return "E3001", true
	case "overloaded":
		return "E3002", true
	case "timeout":
		return "E3003", true
	case "conflict":
		return "E4001", true
	case "cancelled":
		return "E4002", true
	default:
		return "", false
	}
}

func upperOrPOST(m string) string {
	if m == "" {
		return "POST"
	}
	if m == "get" || m == "GET" {
		return "GET"
	}
	if m == "delete" || m == "DELETE" {
		return "DELETE"
	}
	return "POST"
}

// StreamingDecoderFormat returns the decoder format from manifest for ARCH-002.
// Supports "openai_sse", "anthropic_sse", "sse" (defaults to openai_sse).
func StreamingDecoderFormat(m any) string {
	var format, strategy string
	switch v := m.(type) {
	case *V1Manifest:
		if v.Streaming != nil && v.Streaming.Decoder != nil {
			format, strategy = v.Streaming.Decoder.Format, v.Streaming.Decoder.Strategy
		}
	case *V2Manifest:
		if v.Streaming != nil && v.Streaming.Decoder != nil {
			format, strategy = v.Streaming.Decoder.Format, v.Streaming.Decoder.Strategy
		}
	default:
		return "openai_sse"
	}
	if format == "anthropic_sse" || strategy == "anthropic_event_stream" {
		return "anthropic_sse"
	}
	if format == "openai_sse" || format == "sse" || strategy == "openai_chat" {
		return "openai_sse"
	}
	return "openai_sse"
}
