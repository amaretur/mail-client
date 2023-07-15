package service

import (
	"github.com/amaretur/mail-client/internal/dto"
)

type UserRepository interface {
	GetById(id uint) (*dto.User, error)
	CreateOrUpdate(username, password string) (uint, error)
	Delete(uid uint) error
}

type UserService struct {
	secretKey	[]byte
	repo		UserRepository
}

func NewUserService(secretKey string, repo UserRepository) *UserService {
	return &UserService{
		secretKey: []byte(secretKey),
		repo: repo,
	}
}

func (u *UserService) SignIn(username, password string) (uint, error) {
	return u.repo.CreateOrUpdate(username, password)
}

func (u *UserService) GetByID(id uint) (*dto.User, error) {
	return u.repo.GetById(id)
}

func (u *UserService) Delete(uid uint) error {
	return u.repo.Delete(uid)
}
