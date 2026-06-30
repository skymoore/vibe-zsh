package client

import (
	"testing"

	"github.com/skymoore/vibe-zsh/internal/config"
	"github.com/teilomillet/gollm/providers"
)

// TestOpenAICompatibleProviderAddsBearerAuth verifies that the openai-compatible
// provider, unlike gollm's vllm provider, sends an Authorization: Bearer header.
func TestOpenAICompatibleProviderAddsBearerAuth(t *testing.T) {
	registerOpenAICompatibleProvider()

	p, err := providers.GetDefaultRegistry().Get(
		config.ProviderOpenAICompatible, "sk-secret-token", "some-model", nil,
	)
	if err != nil {
		t.Fatalf("registry.Get(%q) returned error: %v", config.ProviderOpenAICompatible, err)
	}

	if got := p.Name(); got != config.ProviderOpenAICompatible {
		t.Errorf("Name() = %q, want %q", got, config.ProviderOpenAICompatible)
	}

	headers := p.Headers()
	if got, want := headers["Authorization"], "Bearer sk-secret-token"; got != want {
		t.Errorf("Authorization header = %q, want %q", got, want)
	}
}

// TestVLLMProviderHasNoAuth documents the gollm behavior that motivated the
// openai-compatible provider: the stock vllm provider never authenticates.
func TestVLLMProviderHasNoAuth(t *testing.T) {
	p, err := providers.GetDefaultRegistry().Get("vllm", "sk-secret-token", "some-model", nil)
	if err != nil {
		t.Fatalf("registry.Get(\"vllm\") returned error: %v", err)
	}
	if _, ok := p.Headers()["Authorization"]; ok {
		t.Error("vllm provider unexpectedly set an Authorization header")
	}
}
