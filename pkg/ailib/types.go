// Package ailib provides the main public API for ai-lib-go.
// ai-lib-go 主入口，提供统一的 Chat API 客户端。
package ailib

import (
	"context"
	"io"
)

// MessageRole represents the role of a message.
// 消息角色，遵循 AI-Protocol 标准定义。
type MessageRole string

const (
	RoleSystem    MessageRole = "system"
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
	RoleTool      MessageRole = "tool"
)

// Message represents a chat message.
// 聊天消息，支持文本和多模态内容。
type Message struct {
	Role         MessageRole     `json:"role"`
	Content      MessageContent  `json:"content"`
	Name         string          `json:"name,omitempty"`
	ToolCalls    []ToolCall      `json:"tool_calls,omitempty"`
	ToolCallID   string          `json:"tool_call_id,omitempty"`
	Metadata     map[string]any  `json:"metadata,omitempty"`
}

// MessageContent represents the content of a message.
// 消息内容，可以是文本或多模态内容块。
type MessageContent interface {
	isMessageContent()
}

// TextContent represents plain text content.
type TextContent struct {
	Text string `json:"text"`
}

func (TextContent) isMessageContent() {}

// MultiContent represents multiple content blocks.
type MultiContent []ContentBlock

func (MultiContent) isMessageContent() {}

// ContentBlock represents a single content block in a multimodal message.
// 内容块，支持文本、图片、音频、视频等。
type ContentBlock struct {
	Type     string    `json:"type"` // text, image, audio, video
	Text     string    `json:"text,omitempty"`
	Image    *ImageContent `json:"image,omitempty"`
	Audio    *AudioContent `json:"audio,omitempty"`
	Video    *VideoContent `json:"video,omitempty"`
}

// ImageContent represents image content.
type ImageContent struct {
	URL    string `json:"url,omitempty"`
	Base64 string `json:"base64,omitempty"`
	MimeType string `json:"mime_type,omitempty"` // image/png, image/jpeg, etc.
}

// AudioContent represents audio content.
type AudioContent struct {
	URL      string `json:"url,omitempty"`
	Base64   string `json:"base64,omitempty"`
	MimeType string `json:"mime_type,omitempty"` // audio/mp3, audio/wav, etc.
}

// VideoContent represents video content.
type VideoContent struct {
	URL      string `json:"url,omitempty"`
	Base64   string `json:"base64,omitempty"`
	MimeType string `json:"mime_type,omitempty"` // video/mp4, etc.
}

// ToolCall represents a tool call in a message.
// 工具调用，用于 function calling。
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"` // function
	Function FunctionCall `json:"function"`
}

// FunctionCall represents a function call.
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // JSON-encoded arguments
}

// ToolDefinition represents a tool definition.
// 工具定义，用于向模型声明可用工具。
type ToolDefinition struct {
	Type     string       `json:"type"` // function
	Function FunctionDef  `json:"function"`
}

// FunctionDef represents a function definition.
type FunctionDef struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]any         `json:"parameters"` // JSON Schema
}

// ChatOptions represents options for chat completion.
// Chat 完成选项。
type ChatOptions struct {
	Model            string           `json:"model"`
	Temperature      *float64         `json:"temperature,omitempty"`
	MaxTokens        *int             `json:"max_tokens,omitempty"`
	TopP            *float64         `json:"top_p,omitempty"`
	FrequencyPenalty *float64         `json:"frequency_penalty,omitempty"`
	PresencePenalty *float64          `json:"presence_penalty,omitempty"`
	Stop            []string         `json:"stop,omitempty"`
	Tools           []ToolDefinition `json:"tools,omitempty"`
	ToolChoice      any              `json:"tool_choice,omitempty"`
	ResponseFormat  *ResponseFormat  `json:"response_format,omitempty"`
	Stream          bool             `json:"stream,omitempty"`
	User            string           `json:"user,omitempty"`
	Seed            *int64           `json:"seed,omitempty"`
	Metadata        map[string]any   `json:"metadata,omitempty"`
}

// ResponseFormat represents the response format.
type ResponseFormat struct {
	Type       string `json:"type"` // text, json_object, json_schema
	JSONSchema any    `json:"json_schema,omitempty"`
}

// ChatResponse represents a chat completion response.
// Chat 完成响应。
type ChatResponse struct {
	ID                string          `json:"id"`
	Object            string          `json:"object"`
	Created           int64           `json:"created"`
	Model             string          `json:"model"`
	Choices           []Choice        `json:"choices"`
	Usage            *Usage          `json:"usage,omitempty"`
	SystemFingerprint string          `json:"system_fingerprint,omitempty"`
}

// Choice represents a choice in a response.
type Choice struct {
	Index        int          `json:"index"`
	Message      Message      `json:"message"`
	Delta        *Delta       `json:"delta,omitempty"`
	FinishReason string       `json:"finish_reason"`
	LogProbs     *LogProbs    `json:"logprobs,omitempty"`
}

// Delta represents a streaming delta.
type Delta struct {
	Role      MessageRole `json:"role,omitempty"`
	Content   string      `json:"content,omitempty"`
	ToolCalls []ToolCall  `json:"tool_calls,omitempty"`
}

// LogProbs represents log probabilities.
type LogProbs struct {
	Content []LogProb `json:"content,omitempty"`
}

// LogProb represents a single log probability.
type LogProb struct {
	Token       string  `json:"token"`
	LogProb     float64 `json:"logprob"`
	Bytes      []byte  `json:"bytes,omitempty"`
	TopLogProbs []TopLogProb `json:"top_logprobs,omitempty"`
}

// TopLogProb represents a top log probability.
type TopLogProb struct {
	Token   string  `json:"token"`
	LogProb float64 `json:"logprob"`
	Bytes   []byte  `json:"bytes,omitempty"`
}

// Usage represents token usage.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// StreamingEvent represents a streaming event.
// 流式事件，遵循 AI-Protocol StreamingEvent 标准。
type StreamingEvent struct {
	Type       string      `json:"type"` // PartialContentDelta, ToolCallStarted, etc.
	Delta      string      `json:"delta,omitempty"`
	ToolCall   *ToolCall   `json:"tool_call,omitempty"`
	Error      *EventError `json:"error,omitempty"`
	FinishReason string    `json:"finish_reason,omitempty"`
}

// EventError represents an error in a streaming event.
type EventError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Stream is an iterator for streaming responses.
type Stream interface {
	// Next advances to the next event. Returns false when done.
	Next() bool
	// Event returns the current event.
	Event() StreamingEvent
	// Err returns any error encountered.
	Err() error
	// Close closes the stream.
	Close() error
}

// Client is the main interface for AI interactions.
// AI 客户端接口。
type Client interface {
	// Chat performs a synchronous chat completion.
	Chat(ctx context.Context, messages []Message, opts *ChatOptions) (*ChatResponse, error)
	// ChatStream performs a streaming chat completion.
	ChatStream(ctx context.Context, messages []Message, opts *ChatOptions) (Stream, error)
	// Close releases resources.
	Close() error
}

// ClientBuilder builds Client instances.
type ClientBuilder struct {
	protocolPath string
	protocolURL  string
	protocolData []byte
	apiKey       string
	baseURL      string
	headers      map[string]string
	timeout      int
	maxRetries   int
}

// NewClientBuilder creates a new client builder.
func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{
		headers: make(map[string]string),
	}
}

// WithProtocolPath sets the protocol file path.
func (b *ClientBuilder) WithProtocolPath(path string) *ClientBuilder {
	b.protocolPath = path
	return b
}

// WithProtocolURL sets the protocol URL.
func (b *ClientBuilder) WithProtocolURL(url string) *ClientBuilder {
	b.protocolURL = url
	return b
}

// WithProtocolData sets the protocol data directly.
func (b *ClientBuilder) WithProtocolData(data []byte) *ClientBuilder {
	b.protocolData = data
	return b
}

// WithAPIKey sets the API key.
func (b *ClientBuilder) WithAPIKey(key string) *ClientBuilder {
	b.apiKey = key
	return b
}

// WithBaseURL sets the base URL.
func (b *ClientBuilder) WithBaseURL(url string) *ClientBuilder {
	b.baseURL = url
	return b
}

// WithHeader adds a header.
func (b *ClientBuilder) WithHeader(key, value string) *ClientBuilder {
	b.headers[key] = value
	return b
}

// WithTimeout sets the timeout in seconds.
func (b *ClientBuilder) WithTimeout(seconds int) *ClientBuilder {
	b.timeout = seconds
	return b
}

// WithMaxRetries sets the maximum retries.
func (b *ClientBuilder) WithMaxRetries(retries int) *ClientBuilder {
	b.maxRetries = retries
	return b
}

// Build builds the client.
func (b *ClientBuilder) Build() (Client, error) {
	return newClient(b)
}

// StreamReader is a helper for reading streams.
type StreamReader struct {
	reader  io.ReadCloser
	decoder StreamDecoder
	event   StreamingEvent
	err     error
	done    bool
}

// StreamDecoder decodes streaming events.
type StreamDecoder interface {
	Decode(data []byte) ([]StreamingEvent, error)
}

func (sr *StreamReader) Next() bool {
	if sr.done {
		return false
	}
	// Implementation details in client package
	return false
}

func (sr *StreamReader) Event() StreamingEvent {
	return sr.event
}

func (sr *StreamReader) Err() error {
	return sr.err
}

func (sr *StreamReader) Close() error {
	return sr.reader.Close()
}

// NewTextContent creates a text content.
func NewTextContent(text string) MessageContent {
	return TextContent{Text: text}
}

// NewMultiContent creates a multi-content.
func NewMultiContent(blocks ...ContentBlock) MessageContent {
	return MultiContent(blocks)
}

// NewTextBlock creates a text block.
func NewTextBlock(text string) ContentBlock {
	return ContentBlock{
		Type: "text",
		Text: text,
	}
}

// NewImageBlock creates an image block.
func NewImageBlock(url string) ContentBlock {
	return ContentBlock{
		Type: "image",
		Image: &ImageContent{URL: url},
	}
}

// NewImageBlockBase64 creates an image block from base64.
func NewImageBlockBase64(base64, mimeType string) ContentBlock {
	return ContentBlock{
		Type: "image",
		Image: &ImageContent{
			Base64:   base64,
			MimeType: mimeType,
		},
	}
}

// NewAudioBlock creates an audio block.
func NewAudioBlock(url string) ContentBlock {
	return ContentBlock{
		Type: "audio",
		Audio: &AudioContent{URL: url},
	}
}

// NewVideoBlock creates a video block.
func NewVideoBlock(url string) ContentBlock {
	return ContentBlock{
		Type: "video",
		Video: &VideoContent{URL: url},
	}
}

// UserMessage creates a user message.
func UserMessage(content string) Message {
	return Message{
		Role:    RoleUser,
		Content: NewTextContent(content),
	}
}

// UserMessageMulti creates a user message with multi-content.
func UserMessageMulti(blocks ...ContentBlock) Message {
	return Message{
		Role:    RoleUser,
		Content: NewMultiContent(blocks...),
	}
}

// AssistantMessage creates an assistant message.
func AssistantMessage(content string) Message {
	return Message{
		Role:    RoleAssistant,
		Content: NewTextContent(content),
	}
}

// SystemMessage creates a system message.
func SystemMessage(content string) Message {
	return Message{
		Role:    RoleSystem,
		Content: NewTextContent(content),
	}
}

// ToolResultMessage creates a tool result message.
func ToolResultMessage(toolCallID, content string) Message {
	return Message{
		Role:      RoleTool,
		Content:   NewTextContent(content),
		ToolCallID: toolCallID,
	}
}
