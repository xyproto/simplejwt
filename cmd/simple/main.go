package main

import (
    "fmt"
    "time"

    "github.com/xyproto/simplejwt"
)

func main() {
    // Set the JWT secret
    simplejwt.SetSecret("your-secret-key")

    // Generate a token
    payload := simplejwt.Payload{
        Sub: "1234567890",
        Exp: time.Now().Add(time.Hour).Unix(),
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
