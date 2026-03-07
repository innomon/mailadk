package token

import (
	"crypto"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base32"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// SPAClaims represents the claims in the token issued by the SPA.
type SPAClaims struct {
	LipID string `json:"lip_id"`
	jwt.RegisteredClaims
}

// FinalClaims represents the claims in the final token issued by mailadk.
type FinalClaims struct {
	LipID string `json:"lip_id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type Verifier struct {
	spaPublicKey       crypto.PublicKey
	verifierPrivateKey crypto.PrivateKey
}

func NewVerifier(spaPubKeyPEM, verifierPrivKeyPEM string) (*Verifier, error) {
	v := &Verifier{}

	if spaPubKeyPEM != "" {
		pub, err := parsePublicKey(spaPubKeyPEM)
		if err != nil {
			return nil, fmt.Errorf("failed to parse SPA public key: %w", err)
		}
		v.spaPublicKey = pub
	}

	if verifierPrivKeyPEM != "" {
		priv, err := parsePrivateKey(verifierPrivKeyPEM)
		if err != nil {
			return nil, fmt.Errorf("failed to parse verifier private key: %w", err)
		}
		v.verifierPrivateKey = priv
	}

	return v, nil
}

// DecodeBase32 decodes a base32 encoded string (standard encoding, no padding).
func (v *Verifier) DecodeBase32(encoded string) (string, error) {
	// Standard email bodies might have some whitespace or different padding.
	// We'll try to be robust.
	encoded = strings.TrimSpace(encoded)
	
	// Many SPA implementations might use NoPadding. Let's try to handle both.
	decoder := base32.StdEncoding.WithPadding(base32.NoPadding)
	data, err := decoder.DecodeString(encoded)
	if err != nil {
		// Try with padding if no padding fails
		data, err = base32.StdEncoding.DecodeString(encoded)
		if err != nil {
			return "", fmt.Errorf("failed to decode base32: %w", err)
		}
	}
	return string(data), nil
}

// VerifySPAToken verifies the signature of the token issued by the SPA.
func (v *Verifier) VerifySPAToken(tokenString string) (*SPAClaims, error) {
	if v.spaPublicKey == nil {
		return nil, errors.New("SPA public key not configured")
	}

	token, err := jwt.ParseWithClaims(tokenString, &SPAClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is what we expect (e.g., EdDSA or RS256)
		// We'll support both as mentioned in the spec.
		switch token.Method.(type) {
		case *jwt.SigningMethodEd25519:
			return v.spaPublicKey, nil
		case *jwt.SigningMethodRSA:
			return v.spaPublicKey, nil
		default:
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*SPAClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// SignFinalToken signs a new token containing both the lip_id and the email.
func (v *Verifier) SignFinalToken(lipID, email string) (string, error) {
	if v.verifierPrivateKey == nil {
		return "", errors.New("verifier private key not configured")
	}

	now := time.Now()
	claims := FinalClaims{
		LipID: lipID,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%s:%s", lipID, email),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(1 * time.Hour)), // 1 hour expiry
			Issuer:    "mailadk",
		},
	}

	var method jwt.SigningMethod
	switch v.verifierPrivateKey.(type) {
	case ed25519.PrivateKey:
		method = jwt.SigningMethodEdDSA
	case crypto.Signer: // Common for RSA/ECDSA
		// We'll default to RS256 for RSA keys
		method = jwt.SigningMethodRS256
	default:
		return "", errors.New("unsupported private key type")
	}

	token := jwt.NewWithClaims(method, claims)
	return token.SignedString(v.verifierPrivateKey)
}

func parsePublicKey(pemStr string) (crypto.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}
	return x509.ParsePKIXPublicKey(block.Bytes)
}

func parsePrivateKey(pemStr string) (crypto.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	// Try various private key formats
	if key, err := x509.ParsePKCS8PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	if key, err := x509.ParseECPrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	
	return nil, errors.New("failed to parse private key")
}
