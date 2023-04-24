package main

import (
	"fmt"

	"github.com/xyproto/simplejwt"
)

func main() {
	// Set the secret that is used for generating and validating JWT tokens
	simplejwt.SetSecret("hunter1")

	// Generate a token by passing in a subject and for how many seconds the token should last
	token := simplejwt.SimpleGenerate("bob@zombo.com", 3600)
	if token == "" {
		fmt.Println("Failed to generate token")
		return
	}
	fmt.Printf("Generated token: %s\n", token)

	// Validate the token
	decodedSubject := simplejwt.SimpleValidate(token)
	if decodedSubject == "" {
		fmt.Println("Failed to validate token")
		return
	}
	fmt.Printf("Decoded payload, got subject: %s\n", decodedSubject)
}
