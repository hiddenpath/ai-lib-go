# ai-lib-go

`ai-lib-go` 是 ai-lib 生态中的官方 Go 运行时实现。

该项目遵循 `ai-protocol` 协议约束，并与 `ai-lib-rust`、`ai-lib-python`、`ai-lib-ts` 保持跨运行时一致性。

## 文档语言

- English: `README.md`
- 中文: `README_CN.md`（本文件）

## 生态定位

- 协议规范：`ai-protocol`
- 运行时矩阵：`ai-lib-rust` / `ai-lib-python` / `ai-lib-ts` / `ai-lib-go`
- 配套项目：`ai-protocol-mock`、`spiderswitch`
- 站点对齐：`ailib.info`

## 核心约束

- 协议优先（`ARCH-001`）
- 算子管线一致（`ARCH-002`）
- 跨运行时行为一致（`ARCH-003`）
- 合规测试优先（`TEST-001`）

## 目录结构

- `internal/protocol`：manifest 契约、加载与能力元信息
- `internal/stream`：SSE 解码
- `internal/resilience`：重试与退避策略
- `pkg/ailib`：对外 SDK API（Client、Builder、能力模型、Fallback 策略）
- `tests/compliance`：基于共享夹具的合规执行器

## 能力覆盖

- 基础能力：`chat`、`chat_stream`、`embeddings`
- 扩展能力：`batch`、`stt`、`tts`、`rerank`
- 高级能力：`mcp`、`computer_use`、`reasoning`、`video`

## 运行时保证

- 错误分类：`HTTP 状态码 + provider code/type` 联合判定
- 高级能力门禁：未声明能力快速失败（`E1005`）
- 请求预检：端点与必要字段校验
- 回退策略：健康度跟踪 + 熔断（circuit breaker）

## 合规覆盖

运行时直接消费 `ai-protocol/tests/compliance/cases`，包括：

- `01-protocol-loading`
- `02-error-classification`
- `03-message-building`
- `04-streaming`
- `05-request-building`
- `06-resilience`
- `07-advanced-capabilities`（门禁、端点、回退、mock 行为断言）

## 快速开始

```bash
go test ./...
```

使用内联协议构建客户端：

```go
client, err := ailib.NewClientBuilder().
  WithProtocolData([]byte(manifestYAML)).
  WithAPIKey(os.Getenv("API_KEY")).
  Build()
```

## 发布约定

- 采用 semver，并按 ai-lib 矩阵发布节奏推进。
- 运行时发布需与 `ai-protocol` 兼容窗口对齐。
- `ailib.info` 站点对齐在运行时发布收敛后执行。
