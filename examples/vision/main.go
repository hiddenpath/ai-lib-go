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
		ailib.UserMessageMulti(
			ailib.NewTextBlock("What's in this image?"),
			ailib.NewImageBlock("https://upload.wikimedia.org/wikipedia/commons/thumb/d/dd/Gfp-wisconsin-madison-the-nature-boardwalk.jpg/2560px-Gfp-wisconsin-madison-the-nature-boardwalk.jpg"),
		),
	}

	opts := &ailib.ChatOptions{
		Model:    "gpt-4o-mini",
		MaxTokens: intPtr(300),
	}

	response, err := client.Chat(ctx, messages, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response:", response.Choices[0].Message.Content)
}

func intPtr(i int) *int {
	return &i
}
