package libgbust

import (
	"net/url"
)

// Config represents the command line arguments passed to gbust
type Config struct {
	Cookies    []string
	Goroutines int
	RawURL     string
	URL        *url.URL
	Verbose    bool
	Wordlists  []string
}
