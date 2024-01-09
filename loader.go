package orale

type Loader struct {
	FlagValues         map[string][]any
	EnvironmentValues  map[string][]any
	ConfigurationFiles []*File
}
