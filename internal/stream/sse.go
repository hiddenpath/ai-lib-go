// Package stream implements streaming decoders.
// 流式解码模块，当前提供通用 SSE 解析器。
package stream

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"strings"
)

type Event struct {
	Type         string
	Delta        string
	FinishReason string
	ToolCall     any
	Usage        *EventUsage
}

type EventUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type Decoder struct {
	reader *bufio.Reader
	format string
}

// NewSSEDecoder creates a decoder with default openai_sse format (backward compatible).
func NewSSEDecoder(r io.Reader) *Decoder {
	return NewDecoderWithFormat(r, "openai_sse")
}

// NewDecoderWithFormat creates a decoder with manifest-driven format (ARCH-002).
// Format: "openai_sse" | "anthropic_sse"
func NewDecoderWithFormat(r io.Reader, format string) *Decoder {
	if format == "" {
		format = "openai_sse"
	}
	return &Decoder{reader: bufio.NewReader(r), format: format}
}

func (d *Decoder) Next() (Event, bool, error) {
	for {
		line, err := d.reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return Event{}, false, nil
			}
			return Event{}, false, err
		}
		line = bytes.TrimSpace(line)
		if len(line) == 0 || !bytes.HasPrefix(line, []byte("data:")) {
			continue
		}

		payload := strings.TrimSpace(strings.TrimPrefix(string(line), "data:"))
		if payload == "[DONE]" {
			return Event{}, false, nil
		}

		var raw map[string]any
		if err := json.Unmarshal([]byte(payload), &raw); err != nil {
			// Keep consuming stream; malformed chunk should not kill full session.
			continue
		}
		return d.extractEvent(raw), true, nil
	}
}

func (d *Decoder) extractEvent(raw map[string]any) Event {
	if d.format == "anthropic_sse" {
		return extractEventAnthropic(raw)
	}
	return extractEventOpenAI(raw)
}

func extractEventOpenAI(raw map[string]any) Event {
	ev := Event{Type: "PartialContentDelta"}

	if uRaw, ok := raw["usage"].(map[string]any); ok {
		eu := &EventUsage{}
		if v, ok := uRaw["prompt_tokens"].(float64); ok {
			eu.PromptTokens = int(v)
		}
		if v, ok := uRaw["completion_tokens"].(float64); ok {
			eu.CompletionTokens = int(v)
		}
		if v, ok := uRaw["total_tokens"].(float64); ok {
			eu.TotalTokens = int(v)
		}
		ev.Usage = eu
	}

	choices, ok := raw["choices"].([]any)
	if !ok || len(choices) == 0 {
		return ev
	}
	choice, ok := choices[0].(map[string]any)
	if !ok {
		return ev
	}
	if fr, ok := choice["finish_reason"].(string); ok {
		ev.FinishReason = fr
	}
	delta, ok := choice["delta"].(map[string]any)
	if !ok {
		return ev
	}
	if text, ok := delta["content"].(string); ok {
		ev.Delta = text
	}
	if tools, ok := delta["tool_calls"]; ok {
		ev.ToolCall = tools
		ev.Type = "ToolCallDelta"
	}
	return ev
}

func extractEventAnthropic(raw map[string]any) Event {
	ev := Event{Type: "PartialContentDelta"}
	evType, _ := raw["type"].(string)
	switch evType {
	case "content_block_delta":
		delta, _ := raw["delta"].(map[string]any)
		if dt, _ := delta["type"].(string); dt == "text_delta" {
			if text, ok := delta["text"].(string); ok {
				ev.Delta = text
			}
		} else if dt == "input_json_delta" {
			if pj, ok := delta["partial_json"].(string); ok {
				ev.Delta = pj
				ev.Type = "ToolCallDelta"
				ev.ToolCall = map[string]any{"function": map[string]any{"arguments": pj}}
			}
		}
	case "message_delta":
		delta, _ := raw["delta"].(map[string]any)
		if sr, ok := delta["stop_reason"].(string); ok {
			ev.FinishReason = sr
			ev.Type = "StreamEnd"
		}
	}
	return ev
}
