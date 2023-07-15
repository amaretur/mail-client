package app

import (
	"github.com/amaretur/mail-client/internal/usecase/auth"
	"github.com/amaretur/mail-client/internal/usecase/mail"
)

type Usecase struct {
	Auth *auth.AuthUsecase
	Mail *mail.MailUsecase
}

func NewUsecase(service *Service) *Usecase {
	return &Usecase{

		Auth: auth.NewAuthUsecase(
			service.User, service.Setting, service.Auth),

		Mail: mail.NewMailUsecase(
			service.Keys, service.Setting, 
			service.User, service.Cipher, service.Mail),
	}
}
