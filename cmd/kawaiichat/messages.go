package main

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/xyproto/simplejwt"
)

type Message struct {
	Sender    string
	Content   string
	Timestamp time.Time
}

var messageStore = struct {
	sync.RWMutex
	messages []Message
}{
	messages: make([]Message, 0),
}

// Add new sendMessageHandler
func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	payload, err := simplejwt.Validate(token)
	if err != nil {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	var message Message
	err = json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Sanitize the message content
	message.Content = html.EscapeString(message.Content)

	message.Sender = payload.Subject

	message.Timestamp = time.Now()

	//fmt.Printf("ADDING %+v\n", message)
	messageStore.Lock()
	messageStore.messages = append(messageStore.messages, message)
	messageStore.Unlock()

	w.WriteHeader(http.StatusCreated)
}

func messagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	messageStore.RLock()
	defer messageStore.RUnlock()

	response, err := json.Marshal(messageStore.messages)
	if err != nil {
		http.Error(w, "Error marshalling messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// Add SSE endpoint
func messagesSSEHandler(w http.ResponseWriter, r *http.Request) {
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

	// check if client accepts text/event-stream
	if r.Header.Get("Accept") == "text/event-stream" {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// send SSE every second
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-r.Context().Done():
				return
			case <-ticker.C:
				messageStore.RLock()
				messagesJSON, err := json.Marshal(messageStore.messages)
				messageStore.RUnlock()
				if err != nil {
					fmt.Println("Error marshalling messages:", err)
					continue
				}

				// send SSE event
				event := fmt.Sprintf("data: %s\n\n", messagesJSON)
				if _, err := w.Write([]byte(event)); err != nil {
					fmt.Println("Error sending SSE:", err)
					return
				}
				w.(http.Flusher).Flush()
			}
		}
	} else {
		messageStore.RLock()
		defer messageStore.RUnlock()

		response, err := json.Marshal(messageStore.messages)
		if err != nil {
			http.Error(w, "Error marshalling message list", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}
