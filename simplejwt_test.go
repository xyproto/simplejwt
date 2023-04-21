package simplejwt_test

import (
	"testing"
	"time"

	"github.com/xyproto/simplejwt"
)

func TestSimpleJWT(t *testing.T) {
	// Set the JWT secret, used when generating and validating the JWT tokens
	simplejwt.SetSecret("testsecret")

	// Generate a token
	payload := simplejwt.Payload{
		Subject: "1234567890",
		Expires: time.Now().Add(time.Hour),
	}

	token, err := simplejwt.Generate(payload, nil)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Validate the token
	decodedPayload, err := simplejwt.Validate(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if decodedPayload.Subject != payload.Subject {
		t.Errorf("Expected Subject to be %s, got %s", payload.Subject, decodedPayload.Subject)
	}

	if decodedPayload.Expires.Unix() != payload.Expires.Unix() {
		t.Errorf("Expected Expires to be %v, got %v", payload.Expires, decodedPayload.Expires)
	}
}
