package service

import (
	"log"
	"bytes"
	"errors"
	"strings"
	"encoding/json"
	"encoding/base64"
	"github.com/amaretur/mail-client/pkg/cipher"
	"github.com/amaretur/mail-client/pkg/eds"
	"github.com/amaretur/mail-client/pkg/mails"
)

var (
	ErrNotFoudKeys = errors.New("error not found keys")
	ErrFailedFormatMail = errors.New("error formating mail")

	ErrIncorrectVerifyKey = errors.New("incorrect verify key")
	ErrIncorrectSignKey = errors.New("incorrect sign key")
	ErrIncorrectDecryptKey = errors.New("incorrect decrypt key")
	ErrIncorrectEncryptKey = errors.New("incorrect encrypt key")

	ErrSignNotValid = errors.New("the signature is not valid")
	ErrSignCreate = errors.New("signature creation error")
)

const (
	encMark = "@encrypt"
	signMark = "@sign"
)

type CipherRepository interface {
	GetEncryptAndSignKeys(uid uint, 
		interlocutor string) (string, string, error)
	GetDecryptAndVerifyKeys(uid uint, 
		interlocutor string) (string, string, error)
}

type CipherService struct {
	repo CipherRepository
}

func NewCipherService(repo CipherRepository) *CipherService {
	return &CipherService{
		repo: repo,
	}
}

func (c *CipherService) EncryptMail(
	uid uint, encrypt, sign bool, mail *mails.Mail) (*mails.Mail, error) {

	if !encrypt && !sign {
		return mail, nil
	}

	ek, sk, err := c.repo.GetEncryptAndSignKeys(uid, mail.ToMail)
	if err != nil {
		return nil, ErrNotFoudKeys
	}

	data, err := json.Marshal(mail)
	if err != nil {
		return nil, ErrFailedFormatMail
	}

	var subject string

	if encrypt {
		subject += encMark + " "
		data, err = c.encrypt(data, []byte(ek))
		if err != nil {
			return nil, err
		}
	}

	if sign {
		subject += signMark
		data, err = c.sign(data, []byte(sk))
		if err != nil {
			return nil, err
		}
	}

	return c.formatEncodingMail(mail.FromMail, mail.ToMail, subject, data), nil
}

func (c *CipherService) encrypt(data []byte, keyBytes []byte) ([]byte, error) {

	key, err := cipher.ParseRSAPublicKeyFromPemBytes(keyBytes)
	if err != nil {
		return nil, ErrIncorrectEncryptKey
	}

	out := new(bytes.Buffer)
	err = cipher.Encrypt(key, bytes.NewBuffer(data), out)

	return out.Bytes(), err
}

func (c *CipherService) sign(data []byte, keyBytes []byte) ([]byte, error) {

	key, err := eds.ParseDSAPrivateKeyFromPemBytes(keyBytes)
	if err != nil {
		return nil, ErrIncorrectSignKey
	}

	cert, err := eds.Sign(key, data)
	if err != nil {
		return nil, ErrSignCreate
	}

	return append(data, cert...), nil
}

func (c *CipherService) formatEncodingMail(
	from, to, subject string, data []byte) *mails.Mail {

	return &mails.Mail {
		Subject: subject,
		FromMail: from,
		ToMail: to,
		Files: []mails.File {
			mails.File {
				Name: "mail",
				IsHTML: true,
				Data: base64.StdEncoding.EncodeToString(data),
			},
		},
	}
}

func (c *CipherService) DecryptionMails(
	uid uint, mailList []*mails.Mail) {

	for i, mail := range mailList {

		isEncrypt := strings.Contains(mail.Subject, encMark)
		isSigned := strings.Contains(mail.Subject, signMark)

		if !isEncrypt && !isSigned {
			continue
		}

		dk, vk, err := c.repo.GetDecryptAndVerifyKeys(uid, mail.FromMail)
		if err != nil {
			mail.Error = "no encryption/signature validation keys found"
			continue
		}

		data, err := base64.StdEncoding.DecodeString(mail.GetHTML())
		if err != nil {
			mail.Error = "failed to read the contents of the encryp email"
			continue
		}

		if isSigned {
			data, err = c.verify(data, []byte(vk))
			if err != nil {
				mail.Error = err.Error()
				continue
			}
		}

		if isEncrypt {
			data, err = c.decrypt(data, []byte(dk))
			if err != nil {
				mail.Error = err.Error()
				continue
			}
		}

		newMail := new(mails.Mail)

		err = json.NewDecoder(bytes.NewBuffer(data)).Decode(newMail)
		if err != nil {
			mail.Error = "the contents of the encrypted email are not valid"
			continue
		}

		mailList[i] = newMail
	}
}

func (c *CipherService) verify(data []byte, keyBytes []byte) ([]byte, error) {

	key, err := eds.ParseDSAPublicKeyFromPemBytes(keyBytes)
	if err != nil {
		return data, ErrIncorrectVerifyKey
	}

	cert := data[len(data)-eds.SignSize:]
	data = data[:len(data)-eds.SignSize]

	ok, err := eds.Verify(key, cert, data)
	if err != nil {
		log.Println(err)
		return data, ErrSignNotValid
	}

	if !ok {
		return data, ErrSignNotValid
	}

	return data, nil
}

func (c *CipherService) decrypt(data []byte, keyBytes []byte) ([]byte, error) {

	key, err := cipher.ParseRSAPrivateKeyFromPemBytes(keyBytes)
	if err != nil {
		return data, ErrIncorrectDecryptKey
	}

	out := new(bytes.Buffer)
	err = cipher.Decrypt(key, bytes.NewBuffer(data), out)

	return out.Bytes(), err
}  




