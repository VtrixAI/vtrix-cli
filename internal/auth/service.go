package auth

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"time"

	"github.com/VtrixAI/vtrix-cli/internal/clierrors"
	"github.com/google/uuid"
)

const clientID = "vtrix-cli"

// generateKeyPair generates an ephemeral Ed25519 key pair.
// Returns (publicKeyBase64URL, privateKey, error).
func generateKeyPair() (string, ed25519.PrivateKey, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", nil, err
	}
	pubB64 := base64.RawURLEncoding.EncodeToString(pub)
	return pubB64, priv, nil
}

// buildProof computes Base64URL( Ed25519Sign( SHA256("device_code:timestamp:nonce") ) )
func buildProof(priv ed25519.PrivateKey, deviceCode string, timestamp int64, nonce string) string {
	msg := fmt.Sprintf("%s:%d:%s", deviceCode, timestamp, nonce)
	hash := sha256.Sum256([]byte(msg))
	sig := ed25519.Sign(priv, hash[:])
	return base64.RawURLEncoding.EncodeToString(sig)
}

// Login runs the full device flow and returns (accessToken, refreshToken, apiKey).
func Login(openBrowser func(url string) error) (string, string, string, error) {
	pubKey, privKey, err := generateKeyPair()
	if err != nil {
		return "", "", "", fmt.Errorf("failed to generate key pair: %w", err)
	}

	client := NewClient("")
	dc, err := client.RequestDeviceCode(DeviceCodeRequest{
		ClientID:        clientID,
		ClientPublicKey: pubKey,
	})
	if err != nil {
		return "", "", "", fmt.Errorf("failed to connect to vtrix: %w", err)
	}

	authURL := buildVerificationURL(dc.VerificationURI, dc.UserCode)

	fmt.Printf("\nURL:  %s\n", authURL)
	fmt.Printf("Code: %s\n\n", dc.UserCode)
	if err := openBrowser(authURL); err != nil {
		fmt.Println("(Could not open browser automatically. Please visit the URL above.)")
	} else {
		fmt.Println("Opened your browser automatically. If the page does not load, visit the URL above manually.")
	}

	fmt.Println("Waiting for authorization...")

	interval := time.Duration(dc.Interval) * time.Second
	if interval < time.Second {
		interval = 5 * time.Second
	}
	deadline := time.Now().Add(time.Duration(dc.ExpiresIn) * time.Second)

	for time.Now().Before(deadline) {
		time.Sleep(interval)

		timestamp := time.Now().Unix()
		nonce := uuid.New().String()
		proof := buildProof(privKey, dc.DeviceCode, timestamp, nonce)

		result, err := client.PollToken(TokenRequest{
			DeviceCode: dc.DeviceCode,
			Timestamp:  fmt.Sprintf("%d", timestamp),
			Nonce:      nonce,
			Proof:      proof,
		})
		if err != nil {
			continue
		}

		switch result.Status {
		case "authorized", "":
			if result.AccessToken != "" {
				return result.AccessToken, result.RefreshToken, result.APIKey, nil
			}
		case "expired":
			return "", "", "", &clierrors.CLIError{
				Message: "authorization code expired",
				Hint:    "Run: vtrix auth login",
			}
		case "pending":
		}
	}

	return "", "", "", &clierrors.CLIError{
		Message: "authorization timed out",
		Hint:    "Run: vtrix auth login",
	}
}

// VerifyToken calls /api/v1/auth/me to validate the token.
func VerifyToken(token string) (*MeResponse, error) {
	return NewClient(token).Me()
}

func buildVerificationURL(verificationURI, userCode string) string {
	if verificationURI == "" || userCode == "" {
		return verificationURI
	}

	u, err := url.Parse(verificationURI)
	if err != nil {
		return verificationURI
	}

	q := u.Query()
	q.Set("code", userCode)
	u.RawQuery = q.Encode()
	return u.String()
}
