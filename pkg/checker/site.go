package checker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	// Time allowed to for website to return a reponse.
	timeout = 800 * time.Millisecond
)

// Site represents the Site composition.
type Site struct {
	URL        string
	StatusCode int
	Active     bool
	Valid      bool
	Error      string
}

// Validate a site url and update the `Valid` field.
func (s *Site) Validate() {
	_, err := url.ParseRequestURI(s.URL)
	if s.Valid = (nil == err); !s.Valid {
		s.Error = fmt.Sprintf("Invalid URI: %s", s.URL)
	}
}

// HealthCheck calls the site url and update the `Active` field once a response
// is received based on a specific timeout.
func (s *Site) HealthCheck() {
	client := &http.Client{
		Timeout: time.Duration(timeout),
	}

	resp, err := client.Get(s.URL)
	if s.Active = (nil == err); !s.Active {
		s.Error = err.Error()
		return
	}

	s.StatusCode = resp.StatusCode
}

// Marshal returns JSON encoding of site.
func (s *Site) Marshal() []byte {
	byte, _ := json.Marshal(s)

	return byte
}
