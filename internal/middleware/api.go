package middleware

import (
	"net/http"
)

func TypeJSONMiddleware(h http.Handler) http.Handler {    
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {        
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Headers", "Cookie, Set-Cookie")
		h.ServeHTTP(w, r)    
	})
}
