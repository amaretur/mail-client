package app

import (
	"github.com/amaretur/mail-client/internal/service"

	"github.com/amaretur/mail-client/internal/handler"
	"github.com/amaretur/mail-client/internal/middleware"
)

type Transport struct {
	Http *handler.Handler
}

func NewTransport(u *Usecase, authService *service.AuthService) *Transport {

	amw := middleware.NewAuth(authService)
	cors := middleware.NewCORS(
		"http://localhost:8000",
		"Content-Type, Authorization",
		"POST, GET, PUT, DELETE",
		"POST, GET, PUT, DELETE",
		"false",
		"600",
	)

	h := handler.NewHandler()

	router := h.Router()
	router.Methods("OPTIONS")
	router.Use(middleware.TypeJSONMiddleware)
	router.Use(cors.Middleware)

	h.Register(handler.NewAuthHandler(u.Auth, amw), "/auth")
	h.Register(handler.NewKeyHandler(u.Mail, amw), "/keys")
	h.Register(handler.NewSettingHandler(u.Mail, amw), "/setting")
	h.Register(handler.NewMailHandler(u.Mail, amw), "/mail")

	return &Transport{
		Http: h,
	}
}
