package payloads

// HealthcheckResponse contains server health data
type HealthcheckResponse struct {
	Version    string `json:"version,omitempty"`
	DeployTime string `json:"deployed_at"`
}
