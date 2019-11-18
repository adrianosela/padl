package keymgr

// Manager represents key management operations
type Manager interface {
	PutPriv(string, string) error
	GetPriv(string) (string, error)
	PutPub(string, string) error
	GetPub(string) (string, error)
}
