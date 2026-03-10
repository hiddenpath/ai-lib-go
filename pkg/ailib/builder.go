// Package ailib client builder.
// 客户端构建器，确保协议优先与显式配置。
package ailib

import (
	"fmt"
	"net/http"
	"time"
)

type ClientBuilder struct {
	protocolPath string
	protocolData []byte
	apiKey       string
	baseURL      string
	headers      map[string]string
	timeout      time.Duration
	maxRetries   int
}

func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{
		headers:    map[string]string{},
		timeout:    30 * time.Second,
		maxRetries: 3,
	}
}

func (b *ClientBuilder) WithProtocolPath(path string) *ClientBuilder {
	b.protocolPath = path
	return b
}

func (b *ClientBuilder) WithProtocolData(data []byte) *ClientBuilder {
	b.protocolData = data
	return b
}

func (b *ClientBuilder) WithAPIKey(key string) *ClientBuilder {
	b.apiKey = key
	return b
}

func (b *ClientBuilder) WithBaseURL(url string) *ClientBuilder {
	b.baseURL = url
	return b
}

func (b *ClientBuilder) WithHeader(k, v string) *ClientBuilder {
	b.headers[k] = v
	return b
}

func (b *ClientBuilder) WithTimeout(d time.Duration) *ClientBuilder {
	b.timeout = d
	return b
}

func (b *ClientBuilder) WithMaxRetries(n int) *ClientBuilder {
	b.maxRetries = n
	return b
}

func (b *ClientBuilder) Build() (Client, error) {
	if len(b.protocolData) == 0 && b.protocolPath == "" && b.baseURL == "" {
		return nil, fmt.Errorf("one of protocolData/protocolPath/baseURL must be configured")
	}

	httpClient := &http.Client{Timeout: b.timeout}
	return newClient(b, httpClient)
}
