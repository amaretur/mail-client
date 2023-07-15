package service

import (
	"fmt"
	"time"
	"errors"
	"encoding/json"
	"encoding/base64"

	"github.com/amaretur/mail-client/internal/dto"

	"github.com/amaretur/mail-client/pkg/mails"
	"github.com/amaretur/mail-client/pkg/mails/imap"
	"github.com/amaretur/mail-client/pkg/mails/smtp"
)

var (
	ErrDecodeFolderName = errors.New("error on decoding folder name")
)

type MailService struct {}

func NewMailService() *MailService {
	return &MailService{}
}

func (m *MailService) GetFolders(ic *dto.MailClient) ([]mails.Mailbox, error) {

	i, err := imap.New(ic.Username, ic.Password, ic.Host, int(ic.Port))
	if err != nil {
		return nil, err
	}
	defer i.Logout()

	return i.GetFolders()
}

func (m *MailService) GetMails(
	ic *dto.MailClient, folder string, offset, count uint,
) ([]*mails.Mail, uint32, error) {

	decodedFolder, err := base64.StdEncoding.DecodeString(folder)
	if err != nil {
		return nil, 0, ErrDecodeFolderName
	}

	i, err := imap.New(ic.Username, ic.Password, ic.Host, int(ic.Port))
	if err != nil {
		return nil, 0, err
	}
	defer i.Logout()

	return i.GetMails(string(decodedFolder), uint32(offset), uint32(count))
}

func (m *MailService) SendKeys(
	sc *dto.MailClient, interlocutor string, share *dto.Share) error {

	data, err := json.Marshal(share)
	if err != nil {
		return err
	}

	mail := mails.Mail{
		Subject: "@share",
		FromMail: sc.Username,
		ToMail: interlocutor,
		Files: []mails.File{
			mails.File{
				Name: "keys",
				IsHTML: true,
				Data: string(data),
			},
		},
	}

	return m.Send(sc, &mail)
}

func (m *MailService) Send(sc *dto.MailClient, mail *mails.Mail) error {

	s, err := smtp.New(sc.Username, sc.Password, sc.Host, int(sc.Port))
	if err != nil {
		return err
	}

	m.PrepareMail(mail, sc.Username)

	return s.Send(mail)
}

func (m *MailService) PrepareMail(mail *mails.Mail, username string) {

	mail.FromMail = username

	now := time.Now()

	mail.Date = fmt.Sprintf(
		"%02d.%02d.%d", now.Day(), now.Month(), now.Year())

	mail.Time = fmt.Sprintf(
		"%02d:%02d:%02d", now.Hour(), now.Minute(), now.Second())
}
