package compliance

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	ailib "github.com/hiddenpath/ai-lib-go/pkg/ailib"
)

// TestCase represents a compliance test case.
type TestCase struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Provider    string                 `json:"provider"`
	Model       string                 `json:"model"`
	Input       map[string]interface{} `json:"input"`
	Expected    map[string]interface{} `json:"expected"`
}

// TestSuite represents a compliance test suite.
type TestSuite struct {
	Version    string     `json:"version"`
	Protocol   string     `json:"protocol"`
	TestCases  []TestCase `json:"test_cases"`
}

// LoadTestSuite loads a test suite from a file.
func LoadTestSuite(path string) (*TestSuite, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var suite TestSuite
	if err := json.Unmarshal(data, &suite); err != nil {
		return nil, err
	}

	return &suite, nil
}

// Run runs the compliance tests.
func (s *TestSuite) Run(t *testing.T, client ailib.Client) {
	for _, tc := range s.TestCases {
		t.Run(tc.ID, func(t *testing.T) {
			s.runTestCase(t, client, tc)
		})
	}
}

func (s *TestSuite) runTestCase(t *testing.T, client ailib.Client, tc TestCase) {
	switch tc.Type {
	case "message_building":
		s.testMessageBuilding(t, client, tc)
	case "stream_decode":
		s.testStreamDecode(t, client, tc)
	case "event_mapping":
		s.testEventMapping(t, client, tc)
	case "tool_accumulation":
		s.testToolAccumulation(t, client, tc)
	case "parameter_mapping":
		s.testParameterMapping(t, client, tc)
	case "protocol_loading":
		s.testProtocolLoading(t, client, tc)
	case "error_classification":
		s.testErrorClassification(t, client, tc)
	case "retry_decision":
		s.testRetryDecision(t, client, tc)
	default:
		t.Skipf("Unknown test type: %s", tc.Type)
	}
}

func (s *TestSuite) testMessageBuilding(t *testing.T, client ailib.Client, tc TestCase) {
	ctx := context.Background()
	messages := []ailib.Message{
		ailib.UserMessage("Hello"),
	}

	opts := &ailib.ChatOptions{
		Model: tc.Model,
	}

	_, err := client.Chat(ctx, messages, opts)
	if err != nil {
		t.Errorf("Message building failed: %v", err)
	}
}

func (s *TestSuite) testStreamDecode(t *testing.T, client ailib.Client, tc TestCase) {
	ctx := context.Background()
	messages := []ailib.Message{
		ailib.UserMessage("Hello"),
	}

	opts := &ailib.ChatOptions{
		Model: tc.Model,
	}

	stream, err := client.ChatStream(ctx, messages, opts)
	if err != nil {
		t.Fatalf("Failed to create stream: %v", err)
	}
	defer stream.Close()

	eventCount := 0
	for stream.Next() {
		_ = stream.Event()
		eventCount++
	}

	if stream.Err() != nil {
		t.Errorf("Stream error: %v", stream.Err())
	}

	if eventCount == 0 {
		t.Error("Expected at least one event")
	}
}

func (s *TestSuite) testEventMapping(t *testing.T, client ailib.Client, tc TestCase) {
	t.Skip("Event mapping test requires mock server")
}

func (s *TestSuite) testToolAccumulation(t *testing.T, client ailib.Client, tc TestCase) {
	t.Skip("Tool accumulation test requires mock server")
}

func (s *TestSuite) testParameterMapping(t *testing.T, client ailib.Client, tc TestCase) {
	t.Skip("Parameter mapping test requires mock server")
}

func (s *TestSuite) testProtocolLoading(t *testing.T, client ailib.Client, tc TestCase) {
	t.Skip("Protocol loading test requires manifest files")
}

func (s *TestSuite) testErrorClassification(t *testing.T, client ailib.Client, tc TestCase) {
	t.Skip("Error classification test requires mock server")
}

func (s *TestSuite) testRetryDecision(t *testing.T, client ailib.Client, tc TestCase) {
	t.Skip("Retry decision test requires mock server")
}
