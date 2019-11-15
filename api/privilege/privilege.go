package privilege

// Level represents the level of privilege
// that a Padl user has over a given resource
type Level int

const (
	// PrivilegeLvlReader gives the bearer read-only
	// access to a resource
	PrivilegeLvlReader = Level(0)

	// PrivilegeLvlEditor gives the bearer edit
	// access to a resource, including the ability
	// to add and remove readers
	PrivilegeLvlEditor = Level(1)

	// PrivilegeLvlOwner gives the bearer full
	// access to a resource, including the ability
	// to add and remove other users' access
	PrivilegeLvlOwner = Level(2)
)
