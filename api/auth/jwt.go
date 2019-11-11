package auth

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/adrianosela/padl/lib/keys"
	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

// CustomClaims represents claims we wish to make and verify with JWTs
type CustomClaims struct {
	jwt.StandardClaims
	/* _______StandardClaims:______________
	Audience  string `json:"aud,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	Id        string `json:"jti,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	NotBefore int64  `json:"nbf,omitempty"`
	Subject   string `json:"sub,omitempty"`
	_______________________________________
	*/
	Projects []string `json:"projects"`
}

//NewCustomClaims returns a new CustomClaims object
func NewCustomClaims(sub, aud, iss string, lifetime time.Duration, projects []string) *CustomClaims {
	return &CustomClaims{
		StandardClaims: jwt.StandardClaims{
			Audience:  aud,
			ExpiresAt: time.Now().Add(lifetime).Unix(),
			Id:        uuid.Must(uuid.NewV4()).String(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    iss,
			Subject:   sub,
		},
		Projects: projects,
	}
}

// stdClaimsToCustomClaims populates a CustomClaims struct with a given map of std claims
func stdClaimsToCustomClaims(stdClaims *jwt.MapClaims) (*CustomClaims, error) {
	//marshall the std claims
	stdClaimsBytes, err := json.Marshal(stdClaims)
	if err != nil {
		return nil, err
	}
	//Unmarshal onto a CustomClaims object
	var cc *CustomClaims
	err = json.Unmarshal(stdClaimsBytes, cc)
	if err != nil {
		return nil, err
	}
	return cc, nil
}

func newJWT(claims *CustomClaims, signingMethod jwt.SigningMethod) *jwt.Token {
	return jwt.NewWithClaims(signingMethod, claims)
}

// SignJWT signs a token with the authenticator's key
func (a *Authenticator) SignJWT(tk *jwt.Token) (string, error) {
	return tk.SignedString(a.signer)
}

// ValidateJWT returns the claims within a token as a CustomClaims obect and validates its fields
func (a *Authenticator) ValidateJWT(tkString string) (*CustomClaims, error) {
	var cc CustomClaims
	//parse onto a jwt token object
	keyfunc := func(tk *jwt.Token) (interface{}, error) {
		return a.signer, nil // we use a single key for now
	}
	token, err := jwt.ParseWithClaims(tkString, &cc, keyfunc)
	if err != nil {
		return nil, fmt.Errorf("Could not parse token: %s", err)
	}
	if token == nil || !token.Valid {
		return nil, fmt.Errorf("Token is invalid")
	}
	//We'll only use/check HS512
	if token.Method != jwt.SigningMethodRS512 {
		return nil, fmt.Errorf("Signing Algorithm: %s, not supported", token.Method.Alg())
	}
	// Now to verify individual claims (functions, except groups, inherited from JWT StandardClaims)
	now := time.Now().Unix()
	//Verify text claims
	if !cc.VerifyIssuer(a.iss, true) {
		return nil, fmt.Errorf("Issuer: Expected %s but was %s", a.iss, cc.Issuer)
	}
	if !cc.VerifyAudience(a.aud, true) {
		return nil, fmt.Errorf("Audience: Expected %s but was %s", a.aud, cc.Audience)
	}
	//Verify time claims
	if !cc.VerifyIssuedAt(now, true) {
		return nil, fmt.Errorf("token was used before \"IssuedAt\"")
	}
	if !cc.VerifyExpiresAt(now, true) {
		return nil, fmt.Errorf("token is expired")
	}
	// verify projects here
	return &cc, nil
}

// GenerateJWTForUser generates and signs a token for a given user
func (a *Authenticator) GenerateJWTForUser(email string, projects []string) (string, error) {
	lifetime := time.Duration(time.Hour * 12)
	cc := NewCustomClaims(email, a.aud, a.iss, lifetime, projects)
	tk := newJWT(cc, jwt.SigningMethodRS512)
	tk.Header["kid"] = keys.GetFingerprint(&a.signer.PublicKey)
	signedTk, err := a.SignJWT(tk)
	if err != nil {
		return "", err
	}
	return signedTk, nil
}
