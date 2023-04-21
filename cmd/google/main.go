package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/xyproto/env/v2"
)

func main() {
	if env.Str("GOOGLE_CLIENT_ID") == "" || env.Str("GOOGLE_CLIENT_SECRET") == "" {
		fmt.Println("Please set GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET environment variables.")
		os.Exit(1)
	}

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/generate", generateHandler)
	http.HandleFunc("/protected", protectedHandler)

	http.HandleFunc("/login", googleLoginHandler)

	http.HandleFunc("/callback", googleCallbackHandler)

	fmt.Println("Server running on :4000")
	http.ListenAndServe(":4000", nil)
}
