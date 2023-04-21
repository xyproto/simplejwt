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

// Payload represents the payload of a JWT.
type Payload struct {
	Sub string `json:"sub"`
	Exp int64  `json:"exp"`
}

// JWTHeader represents the header of a JWT.
type JWTHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

var (
	jwtSecret                = []byte("asdfasdf")
	errInvalidTokenFormat    = errors.New("invalid token format")
	errInvalidTokenSignature = errors.New("invalid token signature")
	errInvalidTokenPayload   = errors.New("invalid token payload")
	errTokenExpired          = errors.New("token has expired")
)

// SetSecret sets the secret key used for generating and validating JWT tokens.
func SetSecret(secret string) {
	jwtSecret = []byte(secret)
}

// Generate generates a JWT token with the provided payload and an optional custom header.
func Generate(payload Payload, customHeader *JWTHeader) (string, error) {
	header := JWTHeader{
		Alg: "HS256",
		Typ: "JWT",
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
	mac := hmac.New(sha256.New, jwtSecret)
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

	mac := hmac.New(sha256.New, jwtSecret)
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

	if time.Now().Unix() > payload.Exp {
		return Payload{}, errTokenExpired
	}

	return payload, nil
}
