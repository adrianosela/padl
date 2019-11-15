package payloads

import "errors"

// CreateKeyRequest is the payload for creating a new padl KMS key
type CreateKeyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Bits        int    `json:"bits"`
}

// Validate validates a key creation request
func (ckr *CreateKeyRequest) Validate() error {
	if ckr.Name == "" {
		return errors.New("no key name provided")
	}
	if ckr.Description == "" {
		return errors.New("no description provided")
	}
	if ckr.Bits != 4096 &&
		ckr.Bits != 2048 &&
		ckr.Bits != 1024 &&
		ckr.Bits != 512 {
		return errors.New("invalid key bits, values must be one of { 512, 1024, 2048, 4096 }")
	}
	return nil
}
