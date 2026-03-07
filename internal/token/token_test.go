package token

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base32"
	"encoding/pem"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func generateEd25519KeyPair() (pubPEM, privPEM string, err error) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		return "", "", err
	}

	pubBytes, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return "", "", err
	}
	pubPEM = string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes}))

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return "", "", err
	}
	privPEM = string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}))

	return pubPEM, privPEM, nil
}

func TestVerifier_DecodeBase32(t *testing.T) {
	v := &Verifier{}
	input := "Hello World"
	encoded := base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString([]byte(input))

	decoded, err := v.DecodeBase32(encoded)
	if err != nil {
		t.Fatalf("DecodeBase32 failed: %v", err)
	}
	if decoded != input {
		t.Errorf("Expected %s, got %s", input, decoded)
	}

	// Test with padding
	encodedWithPadding := base32.StdEncoding.EncodeToString([]byte(input))
	decodedWithPadding, err := v.DecodeBase32(encodedWithPadding)
	if err != nil {
		t.Fatalf("DecodeBase32 with padding failed: %v", err)
	}
	if decodedWithPadding != input {
		t.Errorf("Expected %s, got %s", input, decodedWithPadding)
	}
}

func TestVerifier_TokenFlow(t *testing.T) {
	spaPub, spaPrivPEM, err := generateEd25519KeyPair()
	if err != nil {
		t.Fatalf("Failed to generate SPA keys: %v", err)
	}

	verifierPubPEM, verifierPriv, err := generateEd25519KeyPair()
	if err != nil {
		t.Fatalf("Failed to generate Verifier keys: %v", err)
	}
	_ = verifierPubPEM

	// Create Verifier instance
	v, err := NewVerifier(spaPub, verifierPriv)
	if err != nil {
		t.Fatalf("NewVerifier failed: %v", err)
	}

	// 1. Create a mock SPA token
	spaPrivBlock, _ := pem.Decode([]byte(spaPrivPEM))
	spaPrivKey, _ := x509.ParsePKCS8PrivateKey(spaPrivBlock.Bytes)
	
	spaClaims := SPAClaims{
		LipID: "lip123",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	spaToken := jwt.NewWithClaims(jwt.SigningMethodEdDSA, spaClaims)
	spaTokenString, err := spaToken.SignedString(spaPrivKey)
	if err != nil {
		t.Fatalf("Failed to sign SPA token: %v", err)
	}

	// 2. Verify SPA token
	verifiedClaims, err := v.VerifySPAToken(spaTokenString)
	if err != nil {
		t.Fatalf("VerifySPAToken failed: %v", err)
	}
	if verifiedClaims.LipID != "lip123" {
		t.Errorf("Expected LipID lip123, got %s", verifiedClaims.LipID)
	}

	// 3. Sign Final Token
	finalTokenString, err := v.SignFinalToken(verifiedClaims.LipID, "user@example.com")
	if err != nil {
		t.Fatalf("SignFinalToken failed: %v", err)
	}

	// 4. (Optional) Verify Final Token to be sure
	// We'll use a temporary verifier for this
	vFinal, _ := NewVerifier(verifierPubPEM, "")
	finalToken, err := jwt.ParseWithClaims(finalTokenString, &FinalClaims{}, func(token *jwt.Token) (interface{}, error) {
		return vFinal.spaPublicKey, nil // Use verifier's public key as "spa" for verification
	})
	if err != nil {
		t.Fatalf("Final token verification failed: %v", err)
	}
	finalClaims := finalToken.Claims.(*FinalClaims)
	if finalClaims.LipID != "lip123" || finalClaims.Email != "user@example.com" {
		t.Errorf("Final claims mismatch: %+v", finalClaims)
	}
}
