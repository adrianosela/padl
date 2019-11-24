package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/adrianosela/padl/lib/keys"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
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
	// ADD CUSTOM CLAIMS HERE
}

// NewCustomClaims returns a new CustomClaims object
func NewCustomClaims(sub, aud, iss string, lifetime time.Duration) *CustomClaims {
	return &CustomClaims{
		StandardClaims: jwt.StandardClaims{
			Audience:  aud,
			ExpiresAt: time.Now().Add(lifetime).Unix(),
			Id:        uuid.New().String(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    iss,
			Subject:   sub,
		},
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
func (a *Authenticator) ValidateJWT(tkString string, allowedAuds ...string) (*CustomClaims, error) {
	var cc CustomClaims
	//parse onto a jwt token object
	keyfunc := func(tk *jwt.Token) (interface{}, error) {
		return &a.signer.PublicKey, nil // we use a single key for now
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
	// check allowed audiences
	if len(allowedAuds) == 0 {
		if !cc.VerifyAudience(a.aud, true) {
			return nil, fmt.Errorf("Audience: Expected %s but was %s", a.aud, cc.Audience)
		}
	} else {
		var passed bool
		for _, aud := range allowedAuds {
			if passed = cc.VerifyAudience(aud, true); passed {
				break
			}
		}
		if !passed {
			return nil, fmt.Errorf("token audience %s not allowed", cc.Audience)
		}
	}
	//Verify time claims
	if !cc.VerifyIssuedAt(now, true) {
		return nil, fmt.Errorf("token was used before \"IssuedAt\"")
	}
	if !cc.VerifyExpiresAt(now, true) {
		return nil, fmt.Errorf("token is expired")
	}
	return &cc, nil
}

// GenerateJWT generates and signs a token for a given user
func (a *Authenticator) GenerateJWT(email string, aud string) (string, error) {

	var lifetime time.Duration

	if aud == PadlAPIAudience {
		lifetime = time.Duration(time.Hour * 12)
	} else if aud == ServiceAccountAudience {
		lifetime = time.Duration(time.Hour * 24 * 365) // FIXME: valid for a year
	} else {
		return "", errors.New("Audience not recognized")
	}

	cc := NewCustomClaims(email, aud, a.iss, lifetime)

	tk := newJWT(cc, jwt.SigningMethodRS512)
	tk.Header["kid"] = keys.GetFingerprint(&a.signer.PublicKey)
	signedTk, err := a.SignJWT(tk)
	if err != nil {
		return "", fmt.Errorf("could not sign JWT: %s", err)
	}
	return signedTk, nil
}
