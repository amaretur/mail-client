package validation

import (
	"net/mail"

	"github.com/amaretur/mail-client/internal/dto"
	"github.com/amaretur/mail-client/pkg/mails"
)

func IsValidMail(data string) bool {

	_, err := mail.ParseAddress(data)
	return err == nil
}

func IsValidSetting(data *dto.Setting) bool {

	if len(data.SMTPHost) == 0 || len(data.IMAPHost) == 0 {
		return false
	}

	if data.SMTPPort == 0 || data.IMAPPort == 0 {
		return false
	}

	return true
}

func IsValidAuthData(data *dto.AuthData) bool {

	if len(data.Password) == 0 {
		return false
	}

	if !IsValidMail(data.Username) {
		return false
	}

	if !IsValidSetting(&data.Setting) {
		return false
	}

	return true
}

func IsValidMsg(data *mails.Mail) bool {
	return !IsValidMail(data.FromMail) || !IsValidMail(data.ToMail)
}

func IsValidReceive(data *dto.Receive) bool {

	if (len(data.Encrypt) + len(data.Verify)) < 2 {
		return false
	}

	return IsValidMail(data.Interlocutor)
}
