package service

import (
	"github.com/amaretur/mail-client/internal/dto"
)

type SettingRepository interface {

	CreateOrUpdate(data *dto.Setting) error

	GetByUID(uid uint) (*dto.Setting, error)

	GetSMTPSettig(uid uint) (*dto.MailClient, error)
	GetIMAPSettig(uid uint) (*dto.MailClient, error)
}

type SettingService struct {
	repo		SettingRepository
}

func NewSettingService(repo SettingRepository) *SettingService {
	return &SettingService{
		repo: repo,
	}
}

func (s *SettingService) Set(data *dto.Setting) error {
	return s.repo.CreateOrUpdate(data)
}

func (s *SettingService) GetByUID(uid uint) (*dto.Setting, error) {
	return s.repo.GetByUID(uid)
}

func (s *SettingService) GetSMTPSettig(uid uint) (*dto.MailClient, error) {
	return s.repo.GetSMTPSettig(uid)
}

func (s *SettingService) GetIMAPSettig(uid uint) (*dto.MailClient, error) {
	return s.repo.GetIMAPSettig(uid)
}

