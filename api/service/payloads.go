package service

type registrationRequest struct {
	Email  string `json:"email"`
	PubKey string `json:"pub"`
}
