package auth

import "net/http"

// will take in the request and return a response
// http.Request contains the POST data (headers, URL, method, etc.)
// 
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte ("Login endpointt is working."))
}