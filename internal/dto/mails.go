package dto

import (
	"github.com/amaretur/mail-client/pkg/mails"
)

type Mails struct {
	Mails	[]*mails.Mail	`json:"mails"`
	Total	uint32			`json:"total"`
}
