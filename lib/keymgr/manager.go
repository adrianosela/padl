package keymgr

// Manager represents key management operations
type Manager interface {
	PutKey(string, string) error
	GetKey(string) error
}
