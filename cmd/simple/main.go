package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/xyproto/simplejwt"
)

func generateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	payload := simplejwt.Payload{
		Sub: "1234567890",
		Exp: time.Now().Add(time.Hour).Unix(),
	}
	token, err := simplejwt.Generate(payload, nil)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(token))
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header not provided", http.StatusUnauthorized)
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		http.Error(w, "Token not provided", http.StatusUnauthorized)
		return
	}

	_, err := simplejwt.Validate(token)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Access granted to protected data."}`))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Simple JWT Example</title>
	<style>
		body {
			font-family: Arial, sans-serif;
			max-width: 800px;
			margin: 0 auto;
			padding: 1rem;
		}
		pre {
			background-color: #f5f5f5;
			padding: 0.5rem;
			overflow-x: scroll;
		}
	</style>
</head>
<body>
	<h1>Simple JWT Example</h1>
	<p>Use the following curl commands to interact with the server:</p>
	<h2>1. Generate a JWT token</h2>
	<p>Send a POST request to <code>/generate</code> to generate a JWT token:</p>
	<pre>curl -X POST http://localhost:4000/generate</pre>
	<h2>2. Access protected data</h2>
	<p>Send a GET request to <code>/protected</code> with the token in the Authorization header to access protected data:</p>
	<pre>curl -H "Authorization: Bearer &lt;your_token_here&gt;" http://localhost:4000/protected</pre>
	<p>Replace <code>&lt;your_token_here&gt;</code> with the token you received from the previous command.</p>
</body>
</html>
`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/generate", generateHandler)
	http.HandleFunc("/protected", protectedHandler)
	fmt.Println("Server running on :4000")
	http.ListenAndServe(":4000", nil)
}
