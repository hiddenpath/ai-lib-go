# Changelog

## [Unreleased]

### Added

- Full repository reset for Go runtime rearchitecture baseline.
- Protocol-driven manifest models and loader (`internal/protocol`).
- Unified public client surface for chat/stream/embeddings/batch/stt/tts/rerank.
- Advanced capability APIs and client methods for mcp/computer_use/reasoning/video.
- Manifest capability gate for advanced features (undeclared capability returns `E1005`).
- Provider-aware error classification (`error.code/type` + HTTP status) with fallback matrix helpers.
- Transport preflight validations for endpoint path/method and required identifiers/messages.
- Manifest-driven error classification hook (`error_classification`) wired into runtime error parser.
- `FallbackClient` expanded to full `Client` surface with health-aware circuit breaker failover (including stream connect failover).
- Generic SSE stream decoder and retry execution wrapper.
- Baseline tests for chat/stream/embeddings, advanced capability routes, error mapping, and preflight behavior.
- Compliance test runner now executes all required ai-protocol categories directly from fixtures:
  - protocol loading
  - error classification
  - message building
  - stream decode / event mapping / tool accumulation
  - request parameter mapping
  - retry decision
