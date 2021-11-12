package config

// Config holds the tsfinder configuration
type Config struct {
	Verbose             bool
	Color               bool
	Sibling             *[]string
	OnlyText            bool
	ExcludelistFilename string
}
