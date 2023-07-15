package app

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/amaretur/mail-client/internal/repository"
)

type PgConfig struct {
	Host		string
	Port		string
	Username	string
	Password	string
	DBname		string
	SSLmode		string
}

func NewPostgresDB(cfg PgConfig) (*sqlx.DB, error) {

	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBname, cfg.Password, cfg.SSLmode))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

type Repository struct {
	User	*repository.UserRepository
	Keys	*repository.KeyRepository
	Setting	*repository.SettingRepository
}

func NewReposotory(db *sqlx.DB) *Repository {
	return &Repository{
		User: repository.NewUserRepository(db),
		Keys: repository.NewKeyRepository(db),
		Setting: repository.NewSettingRepository(db),
	}
}
