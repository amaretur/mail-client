package middleware

import (
	"net/http"
)

type CORS struct {
	AllowOrigin			string
	AllowHeaders		string
	RequestMethod		string
	AllowCredentials	string
	AllowMethods		string
	MaxAge				string
}

func NewCORS(origin, headers, method, methods, credentials, maxAge string) *CORS {
	return &CORS{
		AllowOrigin: origin,
		AllowHeaders: headers,
		RequestMethod: method,
		AllowMethods: methods,
		AllowCredentials: credentials,
		MaxAge: maxAge,
	}
}

func (c *CORS) Middleware(h http.Handler) http.Handler {    
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {        

		w.Header().Set("Access-Control-Allow-Origin", c.AllowOrigin)    
		w.Header().Set("Access-Control-Allow-Headers", c.AllowHeaders)
		w.Header().Set("Access-Control-Request-Method", c.RequestMethod)
		w.Header().Set("Access-Control-Allow-Methods", c.AllowMethods)
		w.Header().Set("Access-Control-Allow-Credentials", c.AllowCredentials)
		w.Header().Set("Access-Control-Max-Age", c.MaxAge)

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)    
	})
}
