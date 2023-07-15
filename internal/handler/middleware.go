package handler

import (
	"net/http"
)

// here are all interfaces describing all 
// middleware that are used by handlers

type AuthMiddleware interface {
	AuthOnly(h http.Handler) http.Handler
}

