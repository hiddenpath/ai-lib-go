package main

import (
	"context"
	"fmt"
	"log"
	"os"

	ailib "github.com/hiddenpath/ai-lib-go/pkg/ailib"
)

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY not set")
	}

	client, err := ailib.NewClientBuilder().
		WithAPIKey(apiKey).
		WithBaseURL("https://api.openai.com/v1").
		Build()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()
	messages := []ailib.Message{
		ailib.SystemMessage("You are a helpful assistant."),
		ailib.UserMessage("Hello, how are you?"),
	}

	opts := &ailib.ChatOptions{
		Model:       "gpt-4o-mini",
		Temperature: floatPtr(0.7),
		MaxTokens:   intPtr(100),
	}

	response, err := client.Chat(ctx, messages, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response:", response.Choices[0].Message.Content)
	fmt.Println("Usage:", response.Usage)

	stream, err := client.ChatStream(ctx, messages, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer stream.Close()

	fmt.Println("\nStreaming:")
	for stream.Next() {
		event := stream.Event()
		fmt.Print(event.Delta)
	}
	if stream.Err() != nil {
		log.Fatal(stream.Err())
	}
	fmt.Println()
}

func floatPtr(f float64) *float64 {
	return &f
}

func intPtr(i int) *int {
	return &i
}
