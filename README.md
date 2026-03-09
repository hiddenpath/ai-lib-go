# ai-lib-go

# ai-lib-go: Go 语言 AI-Protocol 运行时
# Official Go Runtime for AI-Protocol

**ai-lib-go** is a high-performance, idiomatic Go implementation of the AI-Protocol specification,
providing unified AI model interaction across multiple providers.

ai-lib-go 是 AI-Protocol 规范的官方 Go 运行时，提供统一的多厂商 AI 模型交互接口。

## Core Philosophy

- **Protocol-Driven**: All behavior is configured through protocol manifests, not code
- **Provider-Agnostic**: Unified interface across OpenAI, Anthropic, Google, and others
- **Streaming-First**: Native support for Server-Sent Events (SSE) streaming
- **Idiomatic Go**: Uses Go's native concurrency model (goroutines, channels)

## Features

- [x] V1/V2 Protocol Loading
- [x] Unified Chat API (sync/stream)
- [x] SSE Streaming Decoding
- [x] Error Classification & Retry
- [x] Tool Calling Support
- [x] Multimodal Support (Vision/Audio/Video)
- [x] Embeddings API
- [x] Batch Processing
- [x] STT/TTS
- [x] Reranking

## Installation

```bash
go get github.com/hiddenpath/ai-lib-go@latest
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    ailib "github.com/hiddenpath/ai-lib-go/pkg/ailib"
)

func main() {
    // Create client with protocol
    client, err := ailib.NewClientBuilder().
        WithProtocolPath("path/to/openai.yaml").
        WithAPIKey("your-api-key").
        Build()
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Create messages
    messages := []ailib.Message{
        {Role: ailib.RoleUser, Content: "Hello, how are you?"},
    }

    // Synchronous chat
    ctx := context.Background()
    response, err := client.Chat(ctx, messages, nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(response.Content)

    // Streaming chat
    stream, err := client.ChatStream(ctx, messages, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer stream.Close()

    for stream.Next() {
        event := stream.Event()
        fmt.Print(event.Delta)
    }
}
```

## Project Structure

```
ai-lib-go/
├── api/            # Public API definitions
│   ├── v1/         # V1 protocol types
│   └── v2/         # V2 protocol types
├── internal/       # Internal implementation
│   ├── protocol/   # Protocol loader
│   ├── client/     # HTTP client
│   ├── pipeline/   # Streaming pipeline
│   ├── transport/  # Transport layer
│   ├── errors/     # Error handling
│   └── resilience/ # Retry & rate limiting
├── pkg/            # Public packages
│   └── ailib/      # Main client API
├── cmd/            # CLI tools
│   └── ailib-cli/  # Command-line interface
├── examples/       # Usage examples
└── tests/          # Test suites
    ├── compliance/ # Compliance tests
    └── integration/# Integration tests
```

## Architecture

ai-lib-go follows the AI-Protocol operator pipeline model:

```
decode → select → accumulate → fanout → map
```

Each stage is implemented as a separate component, allowing for flexible
configuration and extensibility.

## Protocol Loading

Load provider manifests from local files or remote URLs:

```go
// From local file
client, _ := ailib.NewClientBuilder().
    WithProtocolPath("./protocols/openai.yaml").
    Build()

// From remote URL
client, _ := ailib.NewClientBuilder().
    WithProtocolURL("https://example.com/protocols/anthropic.yaml").
    Build()

// From embedded data
client, _ := ailib.NewClientBuilder().
    WithProtocolData(manifestData).
    Build()
```

## Error Handling

ai-lib-go provides structured error handling following the AI-Protocol standard:

```go
response, err := client.Chat(ctx, messages, nil)
if err != nil {
    var apiErr *ailib.APIError
    if errors.As(err, &apiErr) {
        switch apiErr.Code {
        case ailib.ErrAuthentication:
            // Handle auth error
        case ailib.ErrRateLimited:
            // Handle rate limit
        case ailib.ErrInvalidRequest:
            // Handle invalid request
        }
    }
}
```

## Compliance

ai-lib-go aims to pass all AI-Protocol compliance tests. Run compliance tests:

```bash
go test ./tests/compliance/... -v
```

## Contributing

Contributions are welcome! Please read our contributing guidelines.

## License

Dual-licensed under MIT OR Apache-2.0. See [LICENSE](LICENSE), [LICENSE-MIT](LICENSE-MIT), and [LICENSE-APACHE](LICENSE-APACHE).

## Related Projects

- [ai-protocol](https://github.com/hiddenpath/ai-protocol) - Protocol specification
- [ai-lib-python](https://github.com/hiddenpath/ai-lib-python) - Python runtime
- [ai-lib-rust](https://github.com/hiddenpath/ai-lib-rust) - Rust runtime
- [ai-lib-ts](https://github.com/hiddenpath/ai-lib-ts) - TypeScript runtime
