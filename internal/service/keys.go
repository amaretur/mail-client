package service

import (
	"github.com/amaretur/mail-client/internal/dto"
	"github.com/amaretur/mail-client/pkg/cipher"
	"github.com/amaretur/mail-client/pkg/eds"
)

type KeyRepository interface {

	GetVerify(uid uint) (*dto.VerifyKeys, error)
	UpdateVerify(uid uint, keys *dto.VerifyKeys) error

	GetDialogKeysList(uid, offset, limit uint) (*dto.DialogKeysList, error)
	GetDialogKeys(id uint) (*dto.DialogKeys, error)

	CreateOrUpdateInterlocutorKeys(uid uint, 
		interlocutor, encrypt, verify string) (uint, error)
	CreateOrUpdateUserKeys(uid uint, 
		interlocutor, decrypt, encrypt string) (uint, error)
}

type KeyService struct {
	repo KeyRepository
}

func NewKeyService(repo KeyRepository) *KeyService {
	return &KeyService{
		repo: repo,
	}
}

func (k *KeyService) VeryfiKeyGen() (*dto.VerifyKeys, error) {

	pub, priv, err := eds.DSAKeyGen()
	if err != nil {
		return nil, err
	}

	pubBytes, err := eds.ExportDSAPublicKeyAsPemBytes(pub)
	if err != nil {
		return nil, err
	}

	privBytes, err := eds.ExportDSAPrivateKeyAsPemBytes(priv)
	if err != nil {
		return nil, err
	}

	return &dto.VerifyKeys{
		Public: string(pubBytes),
		Private: string(privBytes),
	}, nil
}

func (k *KeyService) EncryptKeyGen() (*dto.EncKeys, error) {

	pub, priv, err := cipher.RSAKeyGen()
	if err != nil {
		return nil, err
	}

	pubBytes, err := cipher.ExportRSAPublicKeyAsPemBytes(pub)
	if err != nil {
		return nil, err
	}

	privBytes, err := cipher.ExportRSAPrivateKeyAsPemBytes(priv)
	if err != nil {
		return nil, err
	}

	return &dto.EncKeys{
		Public: string(pubBytes),
		Private: string(privBytes),
	}, nil
}

func (k *KeyService) UpdateVerify(uid uint, keys *dto.VerifyKeys) error {
	return k.repo.UpdateVerify(uid, keys)
}

func (k *KeyService) SetRandomVerify(uid uint) error {

	keys, err := k.VeryfiKeyGen()
	if err != nil {
		return err
	}

	return k.repo.UpdateVerify(uid, keys)
}

func (k *KeyService) GetVerify(uid uint) (*dto.VerifyKeys, error) {
	return k.repo.GetVerify(uid)
}

func (k *KeyService) GetDialogKeysList(uid, 
	offset, limit uint) (*dto.DialogKeysList, error) {

	return k.repo.GetDialogKeysList(uid, offset, limit)
}

func (k *KeyService) GetDialogKeys(id uint) (*dto.DialogKeys, error) {
	return k.repo.GetDialogKeys(id)
}

func (k *KeyService) Receive(uid uint, data *dto.Receive) (uint, error) {
	return k.repo.CreateOrUpdateInterlocutorKeys(
		uid, 
		data.Interlocutor, 
		data.Encrypt, 
		data.Verify,
	)
}

func (k *KeyService) Share(
	uid uint, interlocutor string) (*dto.Share, error) {

	keys, err := k.EncryptKeyGen()

	_, err = k.repo.CreateOrUpdateUserKeys(
		uid, 
		interlocutor, 
		keys.Private, 
		keys.Public,
	)
	if err != nil {
		return nil, err
	}

	verify, err := k.repo.GetVerify(uid)

	return &dto.Share{
		Encrypt: keys.Public,
		Verify: verify.Public,
	}, err
}
