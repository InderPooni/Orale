package orale

type Config struct {
	FlagValues         map[string][]any
	EnvironmentValues  map[string][]any
	ConfigurationFiles []*File
}
