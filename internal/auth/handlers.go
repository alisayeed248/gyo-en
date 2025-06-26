package auth

import "net/http"

// will take in the request and return a response
// http.Request contains the POST data (headers, URL, method, etc.)
//
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}
	
	w.Write([]byte("Login endpoint is working."))
}
