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
}

type Decoder struct {
	reader *bufio.Reader
}

func NewSSEDecoder(r io.Reader) *Decoder {
	return &Decoder{reader: bufio.NewReader(r)}
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
		return extractEvent(raw), true, nil
	}
}

func extractEvent(raw map[string]any) Event {
	ev := Event{Type: "PartialContentDelta"}
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
