package smtp

import (
	"io"
	"log"
	"fmt"
	"bytes"
	"errors"
	"strings"
	"encoding/base64"

	"github.com/amaretur/mail-client/pkg/mails"

	smtp "github.com/emersion/go-smtp"
	sasl "github.com/emersion/go-sasl"
)

type SMTP struct {
	Login		string
	Password	string
	Server		string
	Port		int
	Address		string

	Client		sasl.Client
}

func New(login, password, server string, port int) (*SMTP, error) {
	addr := fmt.Sprintf("%s:%d", server, port)

	client := sasl.NewPlainClient("", login, password)

	return &SMTP {
		Login: login,
		Password: password,
		Server: server,
		Port: port,
		Address: addr,
		Client: client,
	}, nil
}

func (s *SMTP) Send(mail *mails.Mail) error {
	d := s.formatMail(mail)
	dBuffer := bytes.NewBuffer(d)
	to := strings.Split(mail.ToMail, ", ")

	err := smtp.SendMail(s.Address, s.Client, s.Login, to, dBuffer)

	if err != io.EOF && err != nil {
		log.Println("smtp:", err)

		return errors.New(
			fmt.Sprintf("smtp: error on sending mail"),
		)
	}

	return nil
}

func (s *SMTP) formatMail(mail *mails.Mail) []byte {
	const (
		b = "0000000000004c321d05efd002bf"
		ctf = "Content-Transfer-Encoding: base64"
		cs = "charset=\"utf-8\""
	)

	// Заголовок письма
	content := "MIME-Version: 1.0\r\n"
	content += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	content += fmt.Sprintf("From: %s\r\n", mail.FromMail)
	content += fmt.Sprintf("To: %s\r\n", mail.ToMail)
	content += fmt.Sprintf(
		"Content-Type: multipart/mixed; boundary=\"%s\"\r\n", b,
	)

	// HTML сообщение
	html := mail.GetHTML()

	content += fmt.Sprintf("\r\n--%s\r\n", b)
	content += fmt.Sprintf("Content-Type: text/html; %s\r\n", cs)
	content += fmt.Sprintf("%s\r\n", ctf)
	content += fmt.Sprintf("\r\n%s\r\n",
		base64.StdEncoding.EncodeToString([]byte(html)),
	)

	// Файлы
	files := mail.GetFiles()

	if len(files) > 0 {
		for _, file := range files {
			content += fmt.Sprintf("\r\n--%s\r\n", b)
			content += fmt.Sprintf(
				"Content-Type: %s; %s\r\n", file.Type, cs,
			)
			content += fmt.Sprintf("%s\r\n", ctf)
			content += fmt.Sprintf(
				"Content-Disposition: attachment; filename=\"%s\"\r\n",
				file.Name,
			)
			content += fmt.Sprintf("\r\n%s", file.Data)
		}
	}

	content += fmt.Sprintf("\r\n--%s--\r\n", b)

	return []byte(content)
}

