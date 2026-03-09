package ailib_test

import (
	"context"
	"testing"

	ailib "github.com/hiddenpath/ai-lib-go/pkg/ailib"
)

func TestClientBuilder(t *testing.T) {
	builder := ailib.NewClientBuilder().
		WithAPIKey("test-key").
		WithBaseURL("https://api.example.com/v1").
		WithTimeout(30)

	client, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build client: %v", err)
	}
	defer client.Close()
}

func TestMessageTypes(t *testing.T) {
	userMsg := ailib.UserMessage("Hello")
	if userMsg.Role != ailib.RoleUser {
		t.Errorf("Expected role user, got %s", userMsg.Role)
	}

	textContent := userMsg.Content.(ailib.TextContent)
	if textContent.Text != "Hello" {
		t.Errorf("Expected text 'Hello', got '%s'", textContent.Text)
	}

	sysMsg := ailib.SystemMessage("You are helpful")
	if sysMsg.Role != ailib.RoleSystem {
		t.Errorf("Expected role system, got %s", sysMsg.Role)
	}

	asstMsg := ailib.AssistantMessage("Hi there")
	if asstMsg.Role != ailib.RoleAssistant {
		t.Errorf("Expected role assistant, got %s", asstMsg.Role)
	}
}

func TestMultiModalContent(t *testing.T) {
	multiMsg := ailib.UserMessageMulti(
		ailib.NewTextBlock("What's this?"),
		ailib.NewImageBlock("https://example.com/image.png"),
	)

	if multiMsg.Role != ailib.RoleUser {
		t.Errorf("Expected role user, got %s", multiMsg.Role)
	}

	content := multiMsg.Content.(ailib.MultiContent)
	if len(content) != 2 {
		t.Errorf("Expected 2 content blocks, got %d", len(content))
	}

	if content[0].Type != "text" {
		t.Errorf("Expected text block, got %s", content[0].Type)
	}

	if content[1].Type != "image" {
		t.Errorf("Expected image block, got %s", content[1].Type)
	}
}

func TestToolMessage(t *testing.T) {
	toolMsg := ailib.ToolResultMessage("call_123", "result data")
	if toolMsg.Role != ailib.RoleTool {
		t.Errorf("Expected role tool, got %s", toolMsg.Role)
	}
	if toolMsg.ToolCallID != "call_123" {
		t.Errorf("Expected tool_call_id 'call_123', got '%s'", toolMsg.ToolCallID)
	}
}

func TestClientChat(t *testing.T) {
	t.Skip("Requires mock server")

	client, err := ailib.NewClientBuilder().
		WithAPIKey("test").
		WithBaseURL("http://localhost:4010/v1").
		Build()
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()
	messages := []ailib.Message{
		ailib.UserMessage("Hello"),
	}

	opts := &ailib.ChatOptions{
		Model: "gpt-4",
	}

	_, err = client.Chat(ctx, messages, opts)
	if err != nil {
		t.Errorf("Chat failed: %v", err)
	}
}

func TestClientStream(t *testing.T) {
	t.Skip("Requires mock server")

	client, err := ailib.NewClientBuilder().
		WithAPIKey("test").
		WithBaseURL("http://localhost:4010/v1").
		Build()
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()
	messages := []ailib.Message{
		ailib.UserMessage("Hello"),
	}

	opts := &ailib.ChatOptions{
		Model: "gpt-4",
	}

	stream, err := client.ChatStream(ctx, messages, opts)
	if err != nil {
		t.Fatal(err)
	}
	defer stream.Close()

	count := 0
	for stream.Next() {
		_ = stream.Event()
		count++
	}
	if stream.Err() != nil {
		t.Errorf("Stream error: %v", stream.Err())
	}
	t.Logf("Received %d events", count)
}
