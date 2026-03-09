package ailib

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"bytes"
)

// AudioTranscriptionRequest represents a speech-to-text request.
type AudioTranscriptionRequest struct {
	File           string `json:"file"`            // File path or URL
	Model          string `json:"model"`
	Language       string `json:"language,omitempty"`
	Prompt         string `json:"prompt,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"` // json, text, srt, verbose_json, vtt
	Temperature    *float64 `json:"temperature,omitempty"`
}

// AudioTranscriptionResponse represents a transcription response.
type AudioTranscriptionResponse struct {
	Text     string  `json:"text"`
	Duration float64 `json:"duration,omitempty"`
	Language string  `json:"language,omitempty"`
	Words    []WordTiming `json:"words,omitempty"`
}

// WordTiming represents word-level timing.
type WordTiming struct {
	Word  string  `json:"word"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

// AudioTranslationRequest represents a translation request.
type AudioTranslationRequest struct {
	File           string `json:"file"`
	Model          string `json:"model"`
	Prompt         string `json:"prompt,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	Temperature    *float64 `json:"temperature,omitempty"`
}

// STTClient provides speech-to-text API.
type STTClient struct {
	client *client
}

// NewSTTClient creates an STT client.
func NewSTTClient(c Client) (*STTClient, error) {
	ac, ok := c.(*client)
	if !ok {
		return nil, fmt.Errorf("invalid client type")
	}
	return &STTClient{client: ac}, nil
}

// Transcribe transcribes audio to text.
func (s *STTClient) Transcribe(ctx context.Context, req *AudioTranscriptionRequest) (*AudioTranscriptionResponse, error) {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := s.client.baseURL + "/audio/transcriptions"
	httpReq, err := newRequest(ctx, "POST", url, reqBytes)
	if err != nil {
		return nil, err
	}

	s.client.setHeaders(httpReq)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, s.client.parseError(resp)
	}

	var response AudioTranscriptionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// Translate translates audio to English text.
func (s *STTClient) Translate(ctx context.Context, req *AudioTranslationRequest) (*AudioTranscriptionResponse, error) {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := s.client.baseURL + "/audio/translations"
	httpReq, err := newRequest(ctx, "POST", url, reqBytes)
	if err != nil {
		return nil, err
	}

	s.client.setHeaders(httpReq)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, s.client.parseError(resp)
	}

	var response AudioTranscriptionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &response, nil
}

// TTSRequest represents a text-to-speech request.
type TTSRequest struct {
	Model       string  `json:"model"`
	Input       string  `json:"input"`
	Voice       string  `json:"voice"`
	ResponseFormat string `json:"response_format,omitempty"` // mp3, opus, aac, flac, wav, pcm
	Speed       *float64 `json:"speed,omitempty"`
}

// TTSResponse represents a TTS response (audio data).
type TTSResponse struct {
	AudioData []byte
	MimeType  string
}

// TTSClient provides text-to-speech API.
type TTSClient struct {
	client *client
}

// NewTTSClient creates a TTS client.
func NewTTSClient(c Client) (*TTSClient, error) {
	ac, ok := c.(*client)
	if !ok {
		return nil, fmt.Errorf("invalid client type")
	}
	return &TTSClient{client: ac}, nil
}

// Create creates audio from text.
func (t *TTSClient) Create(ctx context.Context, req *TTSRequest) (*TTSResponse, error) {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := t.client.baseURL + "/audio/speech"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	t.client.setHeaders(httpReq)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := t.client.http.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, t.client.parseError(resp)
	}

	audioData := make([]byte, 0)
	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			audioData = append(audioData, buf[:n]...)
		}
		if err != nil {
			break
		}
	}

	mimeType := resp.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "audio/mpeg"
	}

	return &TTSResponse{
		AudioData: audioData,
		MimeType:  mimeType,
	}, nil
}
