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

type NewProjResponse struct {
	ID string `json:"id"`
}

type CreateDeployKeyResponse struct {
	DeployKey string `json:"deployKey"`
}

type AddOwnerResponse struct {
	TeamKey string `json:"teamKey"`
}

type AddEditorResponse struct {
	TeamKey string `json:"teamKey"`
}

type AddReaderResponse struct {
	TeamKey string `json:"teamKey"`
}

type AddSecretResponse struct {
	ID string `json:"id"`
}
