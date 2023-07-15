package app

import (
	"github.com/amaretur/mail-client/internal/service"
)

type ServiceCongif struct {
	Repo			*Repository
	SecretKey		string
	AccessExpires	int
	RefreshExpires	int
}

type Service struct {
	Auth		*service.AuthService
	Setting		*service.SettingService
	User		*service.UserService
	Keys		*service.KeyService
	Cipher		*service.CipherService
	Mail		*service.MailService
}

func NewService(conf *ServiceCongif) *Service {
	return &Service{
		Auth: service.NewAuthService(
			conf.SecretKey, 
			conf.AccessExpires, 
			conf.RefreshExpires,
		),
		Setting: service.NewSettingService(conf.Repo.Setting),
		User: service.NewUserService(conf.SecretKey, conf.Repo.User),
		Keys: service.NewKeyService(conf.Repo.Keys),
		Cipher: service.NewCipherService(conf.Repo.Keys),
		Mail: service.NewMailService(),
	}
}
