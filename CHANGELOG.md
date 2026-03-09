# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.5.0] - 2026-03-09

### Added
- Initial release of ai-lib-go
- V1/V2 Protocol Loader with YAML/JSON support
- Unified Chat API (synchronous and streaming)
- SSE streaming decoder
- Standard error codes (E1001-E9999)
- Retry policy with exponential backoff and jitter
- Rate limiting and circuit breaker
- Tool calling support (function calling)
- Multimodal support (vision/audio/video)
- Embeddings API
- Batch processing API
- Speech-to-text (STT) API
- Text-to-speech (TTS) API
- Reranking API
- Compliance test infrastructure

### Architecture
- gRPC/Cloud style project layout
- Protocol-driven design (ARCH-001)
- Cross-runtime consistency with Python/Rust/TS (ARCH-003)
- Standard library only (net/http, encoding/json, etc.)
- Native Go concurrency (goroutines, channels)

### Documentation
- README with quick start guide
- API documentation in code
- Example programs (basic, tool_calling, vision)
- Dual licensing (MIT OR Apache-2.0)
