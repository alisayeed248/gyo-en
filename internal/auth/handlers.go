package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// will take in the request and return a response
// http.Request contains the POST data (headers, URL, method, etc.)
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// fixes Cors issue
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	// parse the request
	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid JSON"))
		return
	}

	fmt.Printf("Username: %s, Password: %s\n", loginReq.Username, loginReq.Password)
	w.Write([]byte(fmt.Sprintf("Received user: %s", loginReq.Username)))
}
