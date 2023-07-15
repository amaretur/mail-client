package middleware

import (
	"strings"
	"context"
	"net/http"
)

type AuthService interface {
	ParseAccessToken(token string) (uint, error)
}

type Auth struct {
	service AuthService
}

func NewAuth(service AuthService) *Auth {
	return &Auth{
		service: service,
	}
}

func (a *Auth) AuthOnly(h http.Handler) http.Handler {    
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {        

		auth := r.Header.Get("Authorization")

		data := strings.Split(auth, " ") // ["Bearer", "<token>"]
		if len(data) != 2 || strings.Compare(data[0], "Bearer") != 0 {
			NewJSONResponseWhithError(w, http.StatusUnauthorized, 
				"invalid authorization header value")
			return
		}

		uid, err := a.service.ParseAccessToken(data[1])
		if err != nil {
			NewJSONResponseWhithError(w, http.StatusUnauthorized, 
				"invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), "uid", uid)

		h.ServeHTTP(w, r.WithContext(ctx)) 
	})
}
