package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/xyproto/env/v2"
	"github.com/xyproto/simplejwt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig = &oauth2.Config{
		ClientID:     env.Str("GOOGLE_CLIENT_ID"),
		ClientSecret: env.Str("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:4000/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
)

func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func googleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Invalid OAuth2 callback", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/generate?access_token=%s", code), http.StatusTemporaryRedirect)
}

func generateHandler(w http.ResponseWriter, r *http.Request) {
	accessToken := r.URL.Query().Get("access_token")
	tokenInfo, err := googleOauthConfig.Exchange(context.Background(), accessToken)
	if err != nil {
		http.Error(w, "Error exchanging access token", http.StatusInternalServerError)
		return
	}

	client := googleOauthConfig.Client(context.Background(), tokenInfo)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		http.Error(w, "Error getting user info", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading user info", http.StatusInternalServerError)
		return
	}

	var userInfo struct {
		Sub string `json:"sub"`
	}
	err = json.Unmarshal(data, &userInfo)
	if err != nil {
		http.Error(w, "Error parsing user info", http.StatusInternalServerError)
		return
	}

	payload := simplejwt.Payload{
		Sub: userInfo.Sub,
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
