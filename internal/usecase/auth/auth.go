package auth

import (
	"github.com/amaretur/mail-client/internal/dto"
)

type AuthService interface {
	ParseRefreshToken(tokenStr string) (uint, error)
	CreateJwtTokens(id uint) (*dto.JwtTokens, error) 
}

type UserService interface {
	SignIn(username, password string) (uint, error)
	Delete(id uint) error
}

type SettingService interface {
	Set(data *dto.Setting) error
}

type AuthUsecase struct {
	user	UserService
	setting	SettingService
	auth	AuthService
}

func NewAuthUsecase(user UserService, 
	setting SettingService, auth AuthService) *AuthUsecase {

	return &AuthUsecase{
		user: user,
		setting: setting,
		auth: auth,
	}
}

func (a *AuthUsecase) SignIn(data *dto.AuthData) (*dto.JwtTokens, error) {

	uid, err := a.user.SignIn(data.Username, data.Password)
	if err != nil {
		return nil, err
	}

	err = a.setting.Set(&dto.Setting{
		UserId: uid,
		IMAPHost: data.Setting.IMAPHost,
		IMAPPort: data.Setting.IMAPPort,
		SMTPHost: data.Setting.SMTPHost,
		SMTPPort: data.Setting.SMTPPort,
	})
	if err != nil {
		return nil, err
	}

	return a.auth.CreateJwtTokens(uid)
}

func (a *AuthUsecase) Refresh(token string) (*dto.JwtTokens, error) {

	userId, err := a.auth.ParseRefreshToken(token)
	if err != nil {
		return nil, err
	}

	return a.auth.CreateJwtTokens(userId)
}

func (a *AuthUsecase) DeleteAccount(uid uint) error {
	return a.user.Delete(uid)
}
