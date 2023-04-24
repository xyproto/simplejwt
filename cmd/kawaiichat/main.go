package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/xyproto/simplejwt"
)

type User struct {
	Nickname string
	Password string
}

var userStore = struct {
	sync.RWMutex
	users map[string]User
}{
	users: make(map[string]User),
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userStore.Lock()
	defer userStore.Unlock()
	if _, exists := userStore.users[user.Nickname]; exists {
		http.Error(w, "Nickname already exists", http.StatusBadRequest)
		return
	}

	userStore.users[user.Nickname] = user
	w.WriteHeader(http.StatusCreated)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userStore.RLock()
	storedUser, exists := userStore.users[user.Nickname]
	userStore.RUnlock()
	if !exists || storedUser.Password != user.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	payload := simplejwt.Payload{
		Subject: user.Nickname,
		Expires: time.Now().Add(time.Hour),
	}
	token, err := simplejwt.Generate(payload, nil)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(token))
}

func logoutHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("You have been logged out"))
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
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

	userStore.RLock()
	defer userStore.RUnlock()

	nicknames := make([]string, 0, len(userStore.users))
	for _, user := range userStore.users {
		nicknames = append(nicknames, user.Nickname)
	}

	response, err := json.Marshal(nicknames)
	if err != nil {
		http.Error(w, "Error marshalling user list", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func main() {
	simplejwt.SetSecret("your-jwt-secret-goes-here")

	http.HandleFunc("/", fileHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/send", sendMessageHandler)
	http.HandleFunc("/messages", messagesHandler)
	http.HandleFunc("/messages/sse", messagesSSEHandler)

	fmt.Println("Serving http://localhost:8080/")
	fmt.Fprintf(os.Stderr, "%v\n", http.ListenAndServe(":8080", nil))
}
