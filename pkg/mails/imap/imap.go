package imap

import (
	"io"
	"fmt"
//	"log"
	"sort"
	"errors"
	"strings"
	"encoding/base64"

	"github.com/amaretur/mail-client/pkg/mails"

	imap "github.com/emersion/go-imap"
	ic "github.com/emersion/go-imap/client"
	mail "github.com/emersion/go-message/mail"

	// Для чтения любой кодировки
	_ "github.com/emersion/go-message/charset"
)

type IMAP struct {
	Login		string
	Password	string
	Server		string
	Port		int

	Client		*ic.Client
}

func New(login, password, server string, port int) (*IMAP, error) {
	addr := fmt.Sprintf("%s:%d", server, port)

	client, err := ic.DialTLS(addr, nil)
	if err != nil {
		return nil, errors.New("imap: error on creating TLS connection")
	}

	if err := client.Login(login, password); err != nil {
		return nil, errors.New("imap: error on login to account")
	}

	return &IMAP {
		Login: login,
		Password: password,
		Server: server,
		Port: port,
		Client: client,
	}, nil
}

func (i *IMAP) Logout() error {
	err := i.Client.Logout()

	if err != nil {
		return errors.New("imap: can't logout")
	} else {
		return nil
	}
}

func (i *IMAP) GetFolders() ([]mails.Mailbox, error) {
	mboxes := make(chan *imap.MailboxInfo, 8)
	res := make(chan error, 1)

	go func() {
		res <- i.Client.List("", "*", mboxes)
	}()

	var result []mails.Mailbox

	j := 1

	for m := range mboxes {
		// Для совместимости: нельзя получить письма из Gmail без
		// специального префикса [Gmail]/, но отображать этот префикс
		// нет необходимости

		prefix, mbox, found := strings.Cut(m.Name, m.Delimiter)
		delimiter := m.Delimiter

		if !found {
			mbox = m.Name
		}

		if 0 == strings.Compare(mbox, "[Gmail]") {
			continue
		}

		if prefix == mbox {
			prefix = ""
			delimiter = ""
		}

		fullName := prefix + delimiter + mbox

		result = append(result, mails.Mailbox {
			Id: base64.StdEncoding.EncodeToString([]byte(fullName)),
			Name: mbox,
			Prefix: prefix,
			Delimiter: delimiter,
		})

		j += 1
	}

	if err := <-res; err != nil {
		return nil, errors.New("imap: error on getting mailbox list")
	}

	return result, nil
}

func (i *IMAP) GetMails(
	mbox string, offset, count uint32,
) ([]*mails.Mail, uint32, error) {

	// Для получения сообщения из какой-то папки, необходимо
	// её выбрать. Поскольку до этого может быть выбрана какая-то другая,
	// сначала снимается выбор с предыдущей
	i.Client.Unselect()

	// Выбор почтового ящика
	status, err := i.Client.Select(mbox, true)
	if err != nil {
		return nil, 0, errors.New(
			fmt.Sprintf("imap: error on getting mails from %q", mbox),
		)
	}

	// Получение сообщений
	if status.Messages == 0 {
		return make([]*mails.Mail, 0), 0, nil
	}

	if offset >= status.Messages {
		return make([]*mails.Mail, 0), 0, nil
	}

	if count + offset > status.Messages {
		count = status.Messages - offset
	}

	// Последние count сообщений со смещением offset
	from := status.Messages - (offset + count) + 1
	to := from + count - 1

	seq := new(imap.SeqSet)
	seq.AddRange(from, to)

	mailList := make(chan *imap.Message, count)
	res := make(chan error, 1)

	fetchItems := []imap.FetchItem {
		// Информация о сообщении (получатель, отправитель, время, ...)
		imap.FetchEnvelope,
		// Информация о MIME
		imap.FetchRFC822,
	}

	go func() {
		res <- i.Client.Fetch(seq, fetchItems, mailList)
	}()

	var result []*mails.Mail

	for m := range mailList {
		// Поскольку может быть несколько отправителей и получателей,
		// необходимо собирать их в массив и объединять

		var fromName, fromMail, toName, toMail []string

		for _, from := range m.Envelope.From {
			fromName = append(fromName, from.PersonalName)
			fromMail = append(fromMail, from.Address())
		}

		for _, to := range m.Envelope.To {
			toName = append(toName, to.PersonalName)
			toMail = append(toMail, to.Address())
		}

		fN := strings.Join(fromName, ", ")
		fM := strings.Join(fromMail, ", ")
		tN := strings.Join(toName, ", ")
		tM := strings.Join(toMail, ", ")

		// Форматирование даты и времени
		unformattDate := m.Envelope.Date

		date := fmt.Sprintf("%02d.%02d.%d",
			unformattDate.Day(), unformattDate.Month(),
			unformattDate.Year(),
		)

		time := fmt.Sprintf("%02d:%02d:%02d",
			unformattDate.Hour(), unformattDate.Minute(),
			unformattDate.Second(),
		)

		// Получение файлов
		files, err := i.getFiles(fM, date, time, m.Body)
		if err != nil {
			return nil, 0, err
		}

		result = append(result, &mails.Mail {
			Id: m.SeqNum,
			Subject: m.Envelope.Subject,
			FromName: fN,
			FromMail: fM,
			ToName: tN,
			ToMail: tM,
			Date: date,
			Time: time,
			Files: files,
		})
	}

	if err := <-res; err != nil {
		return nil, 0, errors.New("imap: error on getting mails")
	}

	// Сортировка от новых писем к старым
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Id > result[j].Id
	})

	return result, status.Messages, err
}

func (i *IMAP) getFiles(sender, date, time string,
	body map[*imap.BodySectionName]imap.Literal) ([]mails.File, error) {

	// Когда вначале встречается HTML, его забираем как текст письма,
	// остальные вхождения - файлы
	// Необходимо для совместимости с SMTP

	sender = strings.ReplaceAll(sender, ".", "_")
	date = strings.ReplaceAll(date, ".", "-")
	time = strings.ReplaceAll(time, ":", "-")

	res := []mails.File{}

	for _, body := range body {
		r, err := mail.CreateReader(body)
		if err != nil {
			return nil, errors.New("imap: error on creating mail reader")
		}

		for p, err := r.NextPart();
			err != io.EOF;
			p, err = r.NextPart() {

			if err != nil {
				return nil, errors.New("imap: error on reading part of mail")
			}

			switch h := p.Header.(type) {
				case *mail.InlineHeader:
					if strings.Contains(p.Header.Get("Content-Type"),
						"text/html") {

						html, err := io.ReadAll(p.Body)
						if err != nil {
							return nil, errors.New(
								"imap: error on reading html part",
							)
						}

						res = append(res, mails.File {
							Name: fmt.Sprintf("%s-%s-%s.html",
								sender, date, time),
							IsHTML: true,
							Type: "text/html",
							Data: string(html),
						})
					}
				case *mail.AttachmentHeader:
					fileName, err := h.Filename()
					if err != nil {
						fileName = fmt.Sprintf("%s-%s-%s", sender, date, time)
					}

					file, err := io.ReadAll(p.Body)
					if err != nil {
						return nil, errors.New(
							"imap: error on reading file part",
						)
					}

					res = append(res, mails.File {
						Name: fileName,
						IsHTML: false,
						Type: p.Header.Get("Content-Type"),
						Data: base64.StdEncoding.EncodeToString(file),
					})
			}
		}
	}

	return res, nil
}

