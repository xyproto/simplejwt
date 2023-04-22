#!/bin/bash

# Set your domain or IP and port
DOMAIN="localhost:8080"

# Register a user
echo "Registering a new user (Mallory)..."
curl -X POST -H "Content-Type: application/json" -d '{"nickname":"Mallory","password":"mallory_password"}' "$DOMAIN/register"
echo

# Log in as Mallory
echo "Logging in as Mallory..."
TOKEN=$(curl -s -X POST -H "Content-Type: application/json" -d '{"nickname":"Mallory","password":"mallory_password"}' "$DOMAIN/login")
echo "Received token: $TOKEN"
echo

echo "Listening for messages as Mallory..."
curl -H "Accept: text/event-stream" -H "Authorization: Bearer $TOKEN" "$DOMAIN/messages/sse"
echo
