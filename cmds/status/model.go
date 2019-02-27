package status

import "time"

type serviceStatus struct {
	Name           string    `json:"name"`
	URL            string    `json:"url"`
	Up             bool      `json:"up"`
	LastStatusCode int       `json:"lastStatusCode"`
	LastCheck      time.Time `json:"lastCheck"`
	Infra          bool      `json:"infra"`
}
