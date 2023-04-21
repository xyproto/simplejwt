// Package simplejwt provides a simple JWT implementation for generating
// and validating JWT tokens with HMAC SHA256 signatures.
package simplejwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Payload represents the payload of a JWT token
type Payload struct {
	Subject string    `json:"sub"`
	Expires time.Time `json:"exp"`
}

// Header represents the header of a JWT token
type Header struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

var (
	secretKey                = []byte("your-secret-key")
	errInvalidTokenFormat    = errors.New("invalid token format")
	errInvalidTokenSignature = errors.New("invalid token signature")
	errInvalidTokenPayload   = errors.New("invalid token payload")
	errTokenExpired          = errors.New("token has expired")
)

// SetSecret sets the secret key used for generating and validating JWT tokens.
func SetSecret(secret string) {
	secretKey = []byte(secret)
}

// Generate generates a JWT token with the provided payload and an optional custom header.
func Generate(payload Payload, customHeader *Header) (string, error) {
	header := Header{
		Algorithm: "HS256",
		Type:      "JWT",
	}

	if customHeader != nil {
		header = *customHeader
	}

	headerBytes, err := json.Marshal(header)
	if err != nil {
		return "", err
	}

	headerEncoded := base64.URLEncoding.EncodeToString(headerBytes)
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	payloadEncoded := base64.URLEncoding.EncodeToString(payloadBytes)

	token := fmt.Sprintf("%s.%s", headerEncoded, payloadEncoded)
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(token))
	signature := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	token = fmt.Sprintf("%s.%s", token, signature)

	return token, nil
}

// Validate validates a JWT token and returns the decoded payload if the token is valid.
func Validate(token string) (Payload, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return Payload{}, errInvalidTokenFormat
	}

	tokenToSign := fmt.Sprintf("%s.%s", parts[0], parts[1])

	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(tokenToSign))
	expectedSignature := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	if parts[2] != expectedSignature {
		return Payload{}, errInvalidTokenSignature
	}

	payloadBytes, err := base64.URLEncoding.DecodeString(parts[1])
	if err != nil {
		return Payload{}, errInvalidTokenPayload
	}

	var payload Payload
	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return Payload{}, errInvalidTokenPayload
	}

	if time.Now().Unix() > payload.Expires.Unix() {
		return Payload{}, errTokenExpired
	}

	return payload, nil
}
