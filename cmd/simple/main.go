package main

import (
	"fmt"
	"time"

	"github.com/xyproto/simplejwt"
)

func main() {
	// Set the secret that is used for generating and validating JWT tokens
	simplejwt.SetSecret("your-secret-key")

	// Generate a token
	payload := simplejwt.Payload{
		Subject: "1234567890",
		Expires: time.Now().Add(time.Hour),
	}

	token, err := simplejwt.Generate(payload, nil)
	if err != nil {
		fmt.Printf("Failed to generate token: %v\n", err)
		return
	}

	fmt.Printf("Generated token: %s\n", token)

	// Validate the token
	decodedPayload, err := simplejwt.Validate(token)
	if err != nil {
		fmt.Printf("Failed to validate token: %v\n", err)
		return
	}

	fmt.Printf("Decoded payload: %+v\n", decodedPayload)
}
