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
		ailib.UserMessage("What is the weather in Tokyo?"),
	}

	tools := []ailib.ToolDefinition{
		{
			Type: "function",
			Function: ailib.FunctionDef{
				Name:        "get_weather",
				Description: "Get the current weather in a location",
				Parameters: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"location": map[string]any{
							"type":        "string",
							"description": "The city and country, e.g., Tokyo, Japan",
						},
						"unit": map[string]any{
							"type":        "string",
							"enum":        []string{"celsius", "fahrenheit"},
							"description": "The temperature unit",
						},
					},
					"required": []string{"location"},
				},
			},
		},
	}

	opts := &ailib.ChatOptions{
		Model:  "gpt-4o-mini",
		Tools:  tools,
		MaxTokens: intPtr(500),
	}

	response, err := client.Chat(ctx, messages, opts)
	if err != nil {
		log.Fatal(err)
	}

	for _, choice := range response.Choices {
		if len(choice.Message.ToolCalls) > 0 {
			for _, tc := range choice.Message.ToolCalls {
				fmt.Printf("Tool Call: %s\n", tc.Function.Name)
				fmt.Printf("Arguments: %s\n", tc.Function.Arguments)

				result := fmt.Sprintf(`{"temperature": 22, "unit": "celsius", "description": "sunny"}`)

				messages = append(messages, choice.Message)
				messages = append(messages, ailib.ToolResultMessage(tc.ID, result))
			}
		}
	}

	response2, err := client.Chat(ctx, messages, opts)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nFinal Response:", response2.Choices[0].Message.Content)
}

func intPtr(i int) *int {
	return &i
}
