package payloads

// HealthcheckResponse contains server health data
type HealthcheckResponse struct {
	Version    string `json:"version"`
	DeployTime string `json:"deployed_at"`
}

// LoginResponse contains the response to a login request
type LoginResponse struct {
	Token string `json:"token"`
}
