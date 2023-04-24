#!/bin/bash

# Set your domain or IP and port
DOMAIN="localhost:8080"

# Register a user
echo "Registering a new user (Alice)..."
curl -X POST -H "Content-Type: application/json" -d '{"nickname":"Alice","password":"alice_password"}' "$DOMAIN/register"
echo

# Log in as Alice
echo "Logging in as Alice..."
TOKEN=$(curl -s -X POST -H "Content-Type: application/json" -d '{"nickname":"Alice","password":"alice_password"}' "$DOMAIN/login")
echo "Received token: $TOKEN"
echo

# Send a message as Alice
echo "Sending a message as Alice..."
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN" -d '{"content":"😁"}' "$DOMAIN/send"
echo

# Get messages as Alice
echo "Getting messages as Alice..."
curl -H "Authorization: Bearer $TOKEN" "$DOMAIN/messages"
echo

# Register another user (Bob)
echo "Registering a new user (Bob)..."
curl -X POST -H "Content-Type: application/json" -d '{"nickname":"Bob","password":"bob_password"}' "$DOMAIN/register"
echo

# Log in as Bob
echo "Logging in as Bob..."
TOKEN_BOB=$(curl -s -X POST -H "Content-Type: application/json" -d '{"nickname":"Bob","password":"bob_password"}' "$DOMAIN/login")
echo "Received token: $TOKEN_BOB"
echo

# Send a message as Bob
echo "Sending a message as Bob..."
curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer $TOKEN_BOB" -d '{"content":"😎"}' "$DOMAIN/send"
echo

# Get messages as Bob
echo "Getting messages as Bob..."
curl -H "Authorization: Bearer $TOKEN_BOB" "$DOMAIN/messages"
echo
