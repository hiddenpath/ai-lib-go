// Package ailib capability request/response models.
// 能力模型定义，覆盖 embeddings/batch/stt/tts/rerank。
package ailib

type EmbeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type EmbeddingResponse struct {
	Model string          `json:"model"`
	Data  []EmbeddingItem `json:"data"`
}

type EmbeddingItem struct {
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
}

type BatchCreateRequest struct {
	InputFileID      string `json:"input_file_id"`
	Endpoint         string `json:"endpoint"`
	CompletionWindow string `json:"completion_window"`
}

type BatchJob struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at"`
}

type STTRequest struct {
	Model    string `json:"model"`
	File     string `json:"file"`
	Language string `json:"language,omitempty"`
}

type STTResponse struct {
	Text string `json:"text"`
}

type TTSRequest struct {
	Model  string `json:"model"`
	Input  string `json:"input"`
	Voice  string `json:"voice"`
	Format string `json:"response_format,omitempty"`
}

type TTSResponse struct {
	AudioData []byte
	MimeType  string
}

type RerankRequest struct {
	Model     string   `json:"model"`
	Query     string   `json:"query"`
	Documents []string `json:"documents"`
}

type RerankResponse struct {
	Model string       `json:"model"`
	Data  []RerankItem `json:"data"`
}

type RerankItem struct {
	Index          int     `json:"index"`
	RelevanceScore float64 `json:"relevance_score"`
}

type MCPTool struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	InputSchema map[string]any `json:"input_schema,omitempty"`
}

type MCPListToolsResponse struct {
	Tools []MCPTool `json:"tools"`
}

type MCPCallToolRequest struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments,omitempty"`
}

type MCPCallToolResponse struct {
	Content []map[string]any `json:"content,omitempty"`
	IsError bool             `json:"is_error,omitempty"`
}

type ComputerUseAction struct {
	ActionType string         `json:"action_type"`
	Parameters map[string]any `json:"parameters,omitempty"`
}

type ComputerUseRequest struct {
	SessionID string              `json:"session_id"`
	Actions   []ComputerUseAction `json:"actions"`
}

type ComputerUseResponse struct {
	SessionID string           `json:"session_id"`
	Results   []map[string]any `json:"results,omitempty"`
}

type ReasoningRequest struct {
	Model     string   `json:"model"`
	Prompt    string   `json:"prompt"`
	MaxTokens int      `json:"max_tokens,omitempty"`
	Stop      []string `json:"stop,omitempty"`
}

type ReasoningResponse struct {
	Model       string `json:"model"`
	Reasoning   string `json:"reasoning,omitempty"`
	FinalAnswer string `json:"final_answer,omitempty"`
}

type VideoGenerateRequest struct {
	Model      string         `json:"model"`
	Prompt     string         `json:"prompt"`
	DurationS  int            `json:"duration_s,omitempty"`
	Resolution string         `json:"resolution,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
}

type VideoJob struct {
	ID        string `json:"id"`
	Status    string `json:"status"`
	CreatedAt int64  `json:"created_at,omitempty"`
	URL       string `json:"url,omitempty"`
}
