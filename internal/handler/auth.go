package handler

import (
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"

	"github.com/amaretur/mail-client/internal/dto"
	"github.com/amaretur/mail-client/internal/validation"
)

type AuthUsecase interface {
	SignIn(data *dto.AuthData) (*dto.JwtTokens, error)
	Refresh(token string) (*dto.JwtTokens, error)
	DeleteAccount(uid uint) error
}

type AuthHandler struct {
	usecase	AuthUsecase
	amw		AuthMiddleware
}

func NewAuthHandler(usecase AuthUsecase, amw AuthMiddleware) *AuthHandler {
	return &AuthHandler{
		usecase: usecase,
		amw: amw,
	}
}

func (a *AuthHandler) Init(router *mux.Router) {

	router.HandleFunc("/sign-in", a.signIn).Methods("POST")
	router.HandleFunc("/refresh", a.refresh).Methods("POST")

	authonly := router.PathPrefix("").Subrouter()
	authonly.Use(a.amw.AuthOnly)
	authonly.HandleFunc("/account", 
		a.deleteAccount).Methods("DELETE")
}

func (a *AuthHandler) signIn(w http.ResponseWriter, r* http.Request) {

	authData := new(dto.AuthData)

	err := json.NewDecoder(r.Body).Decode(authData)
	if err != nil || !validation.IsValidAuthData(authData) {
		NewJSONResponseWhithError(
			w, http.StatusBadRequest, "invalid json structure")
		return
	}

	tokens, err := a.usecase.SignIn(authData)
	if err != nil {
		NewJSONResponseWhithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	NewJSONResponse(w, tokens)
}

func (a *AuthHandler) deleteAccount(w http.ResponseWriter, r* http.Request) {

	id := r.Context().Value("uid").(uint)

	if err := a.usecase.DeleteAccount(id); err != nil {
		NewJSONResponseWhithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *AuthHandler) refresh(w http.ResponseWriter, r* http.Request) {

	var data struct {
		Token	string	`json:"token" binding:"required"`
	} 

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil || len(data.Token) == 0 {
		NewJSONResponseWhithError(
			w, http.StatusBadRequest, "invalid json structure")
		return
	}

	tokens, err := a.usecase.Refresh(data.Token)
	if err != nil {
		NewJSONResponseWhithError(w, 
			http.StatusInternalServerError, "oops...")
		return
	}

	NewJSONResponse(w, tokens)
}
