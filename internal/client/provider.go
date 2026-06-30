package client

import (
	"sync"

	"github.com/skymoore/vibe-zsh/internal/config"
	"github.com/teilomillet/gollm/providers"
)

// authedOpenAIProvider wraps gollm's vllm provider, which speaks the
// OpenAI-compatible dialect against a user-configured endpoint but
// deliberately omits authentication. We promote every method from the
// embedded provider and override only Headers to inject the Bearer token.
type authedOpenAIProvider struct {
	providers.Provider
	apiKey string
}

// Name identifies this provider in logs and errors. The embedded vllm
// provider would otherwise report "vllm", which is misleading here.
func (p *authedOpenAIProvider) Name() string {
	return config.ProviderOpenAICompatible
}

// Headers returns the embedded provider's headers plus an Authorization
// header when an API key is configured.
func (p *authedOpenAIProvider) Headers() map[string]string {
	headers := p.Provider.Headers()
	if headers == nil {
		headers = make(map[string]string)
	}
	if p.apiKey != "" {
		headers["Authorization"] = "Bearer " + p.apiKey
	}
	return headers
}

var registerOpenAICompatibleOnce sync.Once

// registerOpenAICompatibleProvider registers the openai-compatible provider on
// gollm's default registry (the one gollm.NewLLM consults). It is safe to call
// multiple times; registration happens exactly once.
func registerOpenAICompatibleProvider() {
	registerOpenAICompatibleOnce.Do(func() {
		providers.GetDefaultRegistry().Register(
			config.ProviderOpenAICompatible,
			func(apiKey, model string, extraHeaders map[string]string) providers.Provider {
				base := providers.NewVLLMProvider(apiKey, model, extraHeaders)
				return &authedOpenAIProvider{Provider: base, apiKey: apiKey}
			},
		)
	})
}
