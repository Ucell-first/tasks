package communication

const (
	// AllProviders indicates that something should use all available providers.
	// Intended to be used with at least RegisterHandler function.
	AllProviders = "all_providers"

	// DefaultHeadersDelimiter is a default header delimiter. Mostly used in message brokers providers.
	DefaultHeadersDelimiter = "|%|"
)
