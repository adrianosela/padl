package store

// MockDatabase is an in-memory database mock
type MockDatabase struct {
	// TODO: add fields here
}

// NewMockDatabase is the constructor for MockDatabase
func NewMockDatabase() *MockDatabase {
	mdb := &MockDatabase{}
	return mdb
}
