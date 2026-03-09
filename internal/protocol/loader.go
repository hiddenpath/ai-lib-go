// Package protocol implements protocol loading and validation.
// 协议加载器，负责加载和验证 V1/V2 provider manifests。
package protocol

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"

	v1 "github.com/hiddenpath/ai-lib-go/api/v1"
	v2 "github.com/hiddenpath/ai-lib-go/api/v2"
)

// Loader loads and caches protocol manifests.
// 协议加载器，支持从文件、URL 或嵌入数据加载 manifests。
type Loader struct {
	mu       sync.RWMutex
	cache    map[string]interface{} // v1.ProviderManifest or v2.ProviderManifest
	paths    []string               // Search paths for protocol files
	baseURL  string                 // Base URL for remote loading
}

// LoaderOption configures the Loader.
type LoaderOption func(*Loader)

// NewLoader creates a new protocol loader.
func NewLoader(opts ...LoaderOption) *Loader {
	l := &Loader{
		cache: make(map[string]interface{}),
	}
	for _, opt := range opts {
		opt(l)
	}
	return l
}

// WithPaths sets the search paths for protocol files.
func WithPaths(paths ...string) LoaderOption {
	return func(l *Loader) {
		l.paths = paths
	}
}

// WithBaseURL sets the base URL for remote loading.
func WithBaseURL(baseURL string) LoaderOption {
	return func(l *Loader) {
		l.baseURL = baseURL
	}
}

// LoadProvider loads a provider manifest by ID.
// The version is auto-detected from the manifest content.
func (l *Loader) LoadProvider(ctx context.Context, providerID string) (interface{}, error) {
	// Check cache first
	l.mu.RLock()
	if manifest, ok := l.cache[providerID]; ok {
		l.mu.RUnlock()
		return manifest, nil
	}
	l.mu.RUnlock()

	// Try to load from search paths
	for _, path := range l.paths {
		// Try V1 format first
		v1Path := filepath.Join(path, "v1", "providers", providerID+".yaml")
		if manifest, err := l.loadV1Manifest(ctx, v1Path); err == nil {
			l.mu.Lock()
			l.cache[providerID] = manifest
			l.mu.Unlock()
			return manifest, nil
		}

		// Try V2 format
		v2Path := filepath.Join(path, "v2", "providers", providerID+".yaml")
		if manifest, err := l.loadV2Manifest(ctx, v2Path); err == nil {
			l.mu.Lock()
			l.cache[providerID] = manifest
			l.mu.Unlock()
			return manifest, nil
		}

		// Try JSON variants
		v1JSONPath := filepath.Join(path, "v1", "providers", providerID+".json")
		if manifest, err := l.loadV1Manifest(ctx, v1JSONPath); err == nil {
			l.mu.Lock()
			l.cache[providerID] = manifest
			l.mu.Unlock()
			return manifest, nil
		}

		v2JSONPath := filepath.Join(path, "v2", "providers", providerID+".json")
		if manifest, err := l.loadV2Manifest(ctx, v2JSONPath); err == nil {
			l.mu.Lock()
			l.cache[providerID] = manifest
			l.mu.Unlock()
			return manifest, nil
		}
	}

	return nil, fmt.Errorf("provider %s not found in search paths", providerID)
}

// LoadV1Provider loads a V1 provider manifest.
func (l *Loader) LoadV1Provider(ctx context.Context, providerID string) (*v1.ProviderManifest, error) {
	manifest, err := l.LoadProvider(ctx, providerID)
	if err != nil {
		return nil, err
	}

	if v1Manifest, ok := manifest.(*v1.ProviderManifest); ok {
		return v1Manifest, nil
	}

	return nil, fmt.Errorf("provider %s is not a V1 manifest", providerID)
}

// LoadV2Provider loads a V2 provider manifest.
func (l *Loader) LoadV2Provider(ctx context.Context, providerID string) (*v2.ProviderManifest, error) {
	manifest, err := l.LoadProvider(ctx, providerID)
	if err != nil {
		return nil, err
	}

	if v2Manifest, ok := manifest.(*v2.ProviderManifest); ok {
		return v2Manifest, nil
	}

	return nil, fmt.Errorf("provider %s is not a V2 manifest", providerID)
}

// loadV1Manifest loads a V1 manifest from a file.
func (l *Loader) loadV1Manifest(ctx context.Context, path string) (*v1.ProviderManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	manifest := &v1.ProviderManifest{}
	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		if err := yaml.Unmarshal(data, manifest); err != nil {
			return nil, fmt.Errorf("failed to parse V1 YAML manifest: %w", err)
		}
	} else {
		if err := json.Unmarshal(data, manifest); err != nil {
			return nil, fmt.Errorf("failed to parse V1 JSON manifest: %w", err)
		}
	}

	return manifest, nil
}

// loadV2Manifest loads a V2 manifest from a file.
func (l *Loader) loadV2Manifest(ctx context.Context, path string) (*v2.ProviderManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	manifest := &v2.ProviderManifest{}
	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		if err := yaml.Unmarshal(data, manifest); err != nil {
			return nil, fmt.Errorf("failed to parse V2 YAML manifest: %w", err)
		}
	} else {
		if err := json.Unmarshal(data, manifest); err != nil {
			return nil, fmt.Errorf("failed to parse V2 JSON manifest: %w", err)
		}
	}

	return manifest, nil
}

// LoadFromData loads a manifest from raw data.
// Auto-detects version from the content.
func (l *Loader) LoadFromData(data []byte) (interface{}, error) {
	// Try V2 first (more specific structure)
	var v2Manifest v2.ProviderManifest
	if err := yaml.Unmarshal(data, &v2Manifest); err == nil && v2Manifest.ID != "" {
		if v2Manifest.Core.Endpoint.BaseURL != "" {
			return &v2Manifest, nil
		}
	}

	// Try V1
	var v1Manifest v1.ProviderManifest
	if err := yaml.Unmarshal(data, &v1Manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	if v1Manifest.ID == "" {
		return nil, fmt.Errorf("invalid manifest: missing provider ID")
	}

	return &v1Manifest, nil
}

// ClearCache clears the manifest cache.
func (l *Loader) ClearCache() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.cache = make(map[string]interface{})
}

// ListCached returns the list of cached provider IDs.
func (l *Loader) ListCached() []string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	ids := make([]string, 0, len(l.cache))
	for id := range l.cache {
		ids = append(ids, id)
	}
	return ids
}

// DetectVersion detects the protocol version from a manifest.
func DetectVersion(manifest interface{}) (string, error) {
	switch m := manifest.(type) {
	case *v1.ProviderManifest:
		return "v1", nil
	case *v2.ProviderManifest:
		return "v2", nil
	default:
		return "", fmt.Errorf("unknown manifest type: %T", manifest)
	}
}

// GetEndpoint returns the base URL for a provider.
func GetEndpoint(manifest interface{}) (string, error) {
	switch m := manifest.(type) {
	case *v1.ProviderManifest:
		return m.BaseURL, nil
	case *v2.ProviderManifest:
		return m.Core.Endpoint.BaseURL, nil
	default:
		return "", fmt.Errorf("unknown manifest type: %T", manifest)
	}
}

// GetAuthType returns the authentication type for a provider.
func GetAuthType(manifest interface{}) (string, error) {
	switch m := manifest.(type) {
	case *v1.ProviderManifest:
		return m.Auth.Type, nil
	case *v2.ProviderManifest:
		return m.Core.Auth.Type, nil
	default:
		return "", fmt.Errorf("unknown manifest type: %T", manifest)
	}
}

// GetAPIFormat returns the API format for a V1 provider.
// For V2 providers, returns "v2".
func GetAPIFormat(manifest interface{}) (string, error) {
	switch m := manifest.(type) {
	case *v1.ProviderManifest:
		return m.APIFormat, nil
	case *v2.ProviderManifest:
		return "v2", nil
	default:
		return "", fmt.Errorf("unknown manifest type: %T", manifest)
	}
}
