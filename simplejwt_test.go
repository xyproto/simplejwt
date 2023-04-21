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
		Sub: "1234567890",
		Exp: time.Now().Add(time.Hour).Unix(),
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

	if decodedPayload.Sub != payload.Sub {
		t.Errorf("Expected Sub to be %s, got %s", payload.Sub, decodedPayload.Sub)
	}

	if decodedPayload.Exp != payload.Exp {
		t.Errorf("Expected Exp to be %d, got %d", payload.Exp, decodedPayload.Exp)
	}
}
