// Package implements URL structure and ways to create it by parsing raw URL string
package url

import (
	"fmt"
	"regexp"
	"strings"
)

// URL structure represents a parsed URL string
type URL struct {
	Scheme string
	Host   string
	Path   string
}

// Parse rawurl into a URL structure
func Parse(rawURL string) (*URL, error) {
	r := regexp.MustCompile(`^([a-zA-Z]+)://([^/]+)(/.*)?$`)
	parts := r.FindStringSubmatch(rawURL)

	if len(parts) < 4 {
		return nil, fmt.Errorf(
			"unable to parse url string: %q",
			rawURL,
		)
	}

	if parts[3] == "" {
		parts[3] = "/"
	}

	result := &URL{
		Scheme: parts[1],
		Host:   parts[2],
		Path:   parts[3],
	}

	return result, nil
}

// Port is getter of `Port` value out of `Host` field formatted as Hostname:Port
func (url *URL) Port() string {
	parts := regexp.MustCompile(`\s*:\s*`).Split(url.Host, -1)
	port := ""
	if len(parts) > 1 {
		port = parts[len(parts)-1:][0]
	}
	return port
}

// Hostname is getter of `Hostname` value out of `Host` field formatted as Hostname:Port
func (url *URL) Hostname() string {
	parts := regexp.MustCompile(`\s*:\s*`).Split(url.Host, -1)
	hostname := parts[0]
	if len(parts) > 1 {
		hostname = strings.Join(
			parts[:len(parts)-1],
			":",
		)
	}
	return hostname
}

// String makes URL struct implementing the Stringer interface
func (url *URL) String() string {
	return fmt.Sprintf("%s://%s%s", url.Scheme, url.Host, url.Path)
}
