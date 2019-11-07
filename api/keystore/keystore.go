package keystore

type Keystore interface {
	GetKey(string) (string, error)
}
