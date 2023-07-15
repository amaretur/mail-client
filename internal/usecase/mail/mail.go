package mail

import (
	"errors"
	"github.com/amaretur/mail-client/internal/dto"
	"github.com/amaretur/mail-client/pkg/mails"
)

var (
	ErrGetSetting = errors.New("error getting smtp settings")
	ErrGetMails = errors.New("error receiving emails")
)

type KeysService interface {

	GetDialogKeysList(uid, offset, limit uint) (*dto.DialogKeysList, error)
	GetDialogKeys(id uint) (*dto.DialogKeys, error)

	GetVerify(uid uint) (*dto.VerifyKeys, error)
	UpdateVerify(uid uint, keys *dto.VerifyKeys) error
	SetRandomVerify(uid uint) error

	Receive(uid uint, data *dto.Receive) (uint, error)
	Share(uid uint, interlocutor string) (*dto.Share, error)
}

type SettingService interface {

	GetByUID(uid uint) (*dto.Setting, error)
	Set(data *dto.Setting) error 

	GetSMTPSettig(uid uint) (*dto.MailClient, error)
	GetIMAPSettig(uid uint) (*dto.MailClient, error)
}

type UserService interface {
	GetByID(uid uint) (*dto.User, error)
}

type CipherService interface {
	EncryptMail(uid uint, 
		encrypt, sign bool, mail *mails.Mail) (*mails.Mail, error) 
	DecryptionMails(uid uint, mails []*mails.Mail)
}

type MailService interface {
	GetFolders(ic *dto.MailClient) ([]mails.Mailbox, error)
	GetMails (ic *dto.MailClient, 
		folder string, offset, limit uint) ([]*mails.Mail, uint32, error)
	Send(sc *dto.MailClient, mail *mails.Mail) error
	PrepareMail(mail *mails.Mail, username string)
	SendKeys(sc *dto.MailClient, interlocutor string, share *dto.Share) error
}

type MailUsecase struct {
	keys	KeysService
	setting	SettingService
	user	UserService
	cipher	CipherService
	mail	MailService
}

func NewMailUsecase(keys KeysService, setting SettingService, 
	user UserService, cipher CipherService, mail MailService) *MailUsecase {

	return &MailUsecase{
		keys: keys,
		setting: setting,
		user: user,
		cipher: cipher,
		mail: mail,
	}
}

func (m *MailUsecase) GetKeysList(uid, 
	offset, limit uint) (*dto.DialogKeysList, error) {

	return m.keys.GetDialogKeysList(uid, offset, limit)
}

func (m *MailUsecase) GetKeys(id uint) (*dto.DialogKeys, error) {
	return m.keys.GetDialogKeys(id)
}

func (m *MailUsecase) GetVerify(uid uint) (*dto.VerifyKeys, error) {
	return m.keys.GetVerify(uid)
}

func (m *MailUsecase) UpdateVerify(uid uint, keys *dto.VerifyKeys) error {
	return m.keys.UpdateVerify(uid, keys)
}

func (m *MailUsecase) SetRandomVerify(uid uint) error {
	return m.keys.SetRandomVerify(uid)
}

func (m *MailUsecase) Receive(uid uint, data *dto.Receive) (uint, error) {
	return m.keys.Receive(uid, data)
}

func (m *MailUsecase) Share(uid uint, interlocutor string) error {

	share, err := m.keys.Share(uid, interlocutor)
	if err != nil {
		return err
	}

	sc, err := m.setting.GetSMTPSettig(uid)
	if err != nil {
		return err
	}

	return m.mail.SendKeys(sc, interlocutor, share)
}

func (m *MailUsecase) GetSetting(uid uint) (*dto.Setting, error) {
	return m.setting.GetByUID(uid)
}

func (m *MailUsecase) UpdateSetting(data *dto.Setting) error {
	return m.setting.Set(data)
}

func (m *MailUsecase) GetFolders(uid uint) ([]mails.Mailbox, error) {

	ic, err := m.setting.GetIMAPSettig(uid)
	if err != nil {
		return nil, err
	}

	return m.mail.GetFolders(ic)
}

func (m *MailUsecase) GetMails(
	uid uint, folder string, offset, limit uint,
) ([]*mails.Mail, uint32, error) {

	ic, err := m.setting.GetIMAPSettig(uid)
	if err != nil {
		return nil, 0, ErrGetSetting
	}

	mails, total, err := m.mail.GetMails(ic, folder, offset, limit)
	if err != nil {
		return nil, 0, ErrGetMails
	}

	m.cipher.DecryptionMails(uid, mails)

	return mails, total, nil
}

func (m *MailUsecase) SendMail(
	uid uint, encrypt, sign bool, mail *mails.Mail) error {

	sc, err := m.setting.GetSMTPSettig(uid)
	if err != nil {
		return ErrGetSetting
	}

	m.mail.PrepareMail(mail, sc.Username)

	encryptMail, err := m.cipher.EncryptMail(uid, encrypt, sign, mail)
	if err != nil {
		return err
	}

	return m.mail.Send(sc, encryptMail)
}

