package swarmgo

import (
	"net/http"

	"github.com/wlevene/swarmgo/llm"
)

// ClientConfig represents the configuration for an LLM client
type ClientConfig struct {
	llm                llm.LLM
	Provider           llm.LLMProvider
	AuthToken          string
	BaseURL            string
	OrgID              string
	APIVersion         string
	AssistantVersion   string
	ModelMapperFunc    func(model string) string // replace model to provider-specific deployment name
	HTTPClient         *http.Client
	EmptyMessagesLimit uint
	Options            map[string]interface{} // Additional provider-specific options
}
