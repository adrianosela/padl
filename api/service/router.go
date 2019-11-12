package service

import (
	"log"

	"github.com/adrianosela/padl/api/auth"
	"github.com/adrianosela/padl/api/config"
	"github.com/adrianosela/padl/api/store"
	"github.com/adrianosela/padl/lib/keys"
	"github.com/gorilla/mux"
)

// Service holds the service configuration
// necessary for endpoints to respond to requests
type Service struct {
	Config        *config.Config
	Router        *mux.Router
	Database      store.Database
	Authenticator *auth.Authenticator
}

// NewPadlService returns an HTTP router multiplexer with
// attached handler functions
func NewPadlService(c *config.Config) *Service {
	// db := store.NewMongoDB(c.Database.ConnectionString, c.Database.Name, c.Database.UsersCollectionName)
	db := store.NewMockDatabase()

	// FIXME
	priv, err := keys.DecodePrivKeyPEM([]byte(`
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAtfDaWhn5dHXD6Glh4DExOzLVgzfG7RFhMuKa/cT0+kCAYpXY
yqkCwB9BPND7BGGwF3Cd/JIxOSiTjQmAO2607aXtVMMOYQ7feDrdaXRfG5qs/2/3
y3T8wi7N1emeG6TA5h+IZuelvUqDmBRbmklMwthjvjRBumuLcOqOWhqWI9K7FMnG
LlBPCl203c8a4FAujcaJ6rl4fBeCTbmHXUhoW0EpzIRe6iADgYvOgxx3IDQSOo1v
Fv1UqIZ5KFw3vaUN9Bol5OxtxQOvJ0ac2zizi5f6EtzZNGSIL4fdDvZKZg1RVQE4
uCUS8Q1MkLeiq5ZIrEqqykkK6DnkOfoUwumTtQIDAQABAoIBAH2/qir8KN3FR1Ir
A+rgFRbFW60FsAfKK1PwKw+aQXd1fUamKuBnT+9Zqs1N4zB0FDEfNRKMOFk5IkIo
fuiU65gVRqN+7UFH9kwy4zUvqUx663bg/HMyuD+9+aYPgae5h1mGEdCN6o+aILnL
2EQaxWMmDEo58/PUwNuaQikklwDKSwZjzCgUK3auxoI/yJl/qAzX7GKhezXd9SK4
90B0IK6MB+0otL0SOreF/bsWQksa0dvaPHAd6gbD0mNHtfpRbuK+nlhP05d23KVE
sZDnhrXYl8//EsgS0HhgAHPQriNTofqFblps1GlyiaXBoHQ4Vl+dyjzP0FMQ9BSk
ifICRAECgYEA1I3y857qnP/aDnjgKMpfxyjqsMSvem/sqmgF3G/9pqhTKKq2kjXo
iW3A/yHtTGZYZUZ2Ze0hiu73EaJ5chudHoFbvG9TXIdWreztcmgqfrlox9097F8u
b2bmevPKWFKRHs745lm9l1ckhOmWHVRAs2M9v3Ri3Ly/y+2Wikzy6fECgYEA2yEG
SFHnnqDWUQZms132i1A90cOqg2LtBO0gqrFdYNYIu/5aUCccTfvqVuyZy63Y+t07
4iF7EFJ8Lpqght1mYUXQwzbxKrNDRZzoF8otFueHd22lmd7a/PWR9e7yPewMfVFF
fa8u4b+rLmwEo7EUPBQ3osC6P7JMBueann+pIgUCgYATAkLZALxQoBz7MFozq62X
HRSoDF75HytWLglgJm/TyLfvKh07xDBwoe0hpAIZ1AlRvVR3Vxap2ycjX5lm2Atc
IAt5NaeJ3dyln0u48JHkVWaGgUW5buWzNsuj8UuGTJQH4lCmIR5we22bqVwwcUl8
AYMTLTBuNz8b2Lqe0bTjsQKBgDdrh8o8pMbSyMFfTBQrPJKJbckionpuR6HKU0u4
ZfR6zWS2dKL28Uqr3t2zI0aHJmx0DZQogZZkNjIXO2hAkIcjgCQPPjldczMk9vIl
WPgFAJbs7UgYO+xkM1Eu6KdOju4W4uthpgrETggEm7vGqmZzeoq4EaLQdjf81Xcm
tGD9AoGBAJ23QVr+MwnJ7BadWZe9adK767JG/dfNSqXpPuQGWA1sBRUrh3V0AvnA
LJYJncjwYiEd9tm9v/20eVcobp4Kxwd6AesNOe/+kGnUVZwidgpNEzs+BuMEexNt
4CGk85KpVBr9ng1RuG4B4Z2+gUnGTmXkRMweEuEu49qzf7SP/Yh4
-----END RSA PRIVATE KEY-----`))
	if err != nil {
		log.Fatal("could not load mock key")
	}

	svc := &Service{
		Config:        c,
		Router:        mux.NewRouter(),
		Database:      db,
		Authenticator: auth.NewAuthenticator(db, priv, "api.padl.com", "api"),
	}

	svc.addDebugEndpoints()
	svc.addAuthEndpoints()

	return svc
}
