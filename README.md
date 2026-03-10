# ai-lib-go

Official Go runtime for the ai-lib ecosystem.

`ai-lib-go` implements protocol-driven model access and cross-runtime consistency requirements defined by `ai-protocol`, `ai-lib-constitution`, and `ai-lib-plans`.

## Language

- English: `README.md` (this file)
- Chinese: `README_CN.md`

## Position In The Ecosystem

- Protocol source: `ai-protocol`
- Runtime peers: `ai-lib-rust`, `ai-lib-python`, `ai-lib-ts`
- Mock and routing projects: `ai-protocol-mock`, `spiderswitch`
- Site/docs sink: `ailib.info`

## Design Constraints

- Protocol-first behavior (`ARCH-001`)
- Deterministic operator pipeline (`ARCH-002`)
- Cross-runtime consistency (`ARCH-003`)
- Compliance-first delivery (`TEST-001`)

## Package Layout

- `internal/protocol` — manifest model/loader and capability metadata
- `internal/stream` — SSE decoding primitives
- `internal/resilience` — retry and backoff utilities
- `pkg/ailib` — public SDK API (`Client`, builder, capabilities, fallback policy)
- `tests/compliance` — fixture-driven compliance runner

## Supported Capability Surface

- Base: `chat`, `chat_stream`, `embeddings`
- Extended: `batch`, `stt`, `tts`, `rerank`
- Advanced: `mcp`, `computer_use`, `reasoning`, `video`

## Runtime Guarantees

- Manifest-aware error classification (`status + provider code/type`)
- Capability guard for undeclared advanced features (`E1005`)
- Preflight request validation for endpoint shape and required fields
- Fallback executor with health snapshot and circuit-breaker policy

## Compliance Coverage

The runner consumes `ai-protocol/tests/compliance/cases` directly, including:

- `01-protocol-loading`
- `02-error-classification`
- `03-message-building`
- `04-streaming`
- `05-request-building`
- `06-resilience`
- `07-advanced-capabilities` (guard, endpoint, fallback, provider-mock behavior)

## Quick Start

```bash
go test ./...
```

Create a client from inline protocol data:

```go
client, err := ailib.NewClientBuilder().
  WithProtocolData([]byte(manifestYAML)).
  WithAPIKey(os.Getenv("API_KEY")).
  Build()
```

## Release Policy

- Versioning follows semver and matrix release train.
- Runtime release must align with `ai-protocol` compatibility window.
- Public site sync (`ailib.info`) happens after runtime release closure.
