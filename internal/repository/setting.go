package repository

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/amaretur/mail-client/internal/dto"
)

type SettingRepository struct {
	conn *sqlx.DB
}

func NewSettingRepository(db *sqlx.DB) *SettingRepository {
	return &SettingRepository{
		conn: db,
	}
}

func (s *SettingRepository) GetByUID(uid uint) (*dto.Setting, error) {

	var setting dto.Setting

	query := fmt.Sprintf(
		`SELECT 
			user_id, imap_host, imap_port, smtp_host, smtp_port 
		FROM 
			setting 
		WHERE 
			user_id = %d`, uid,
	)

	err := s.conn.Get(&setting, query)
	if err != nil {
		log.Println(err.Error())
	}

	return &setting, err
}

func (s *SettingRepository) CreateOrUpdate(data *dto.Setting) error {

	query := fmt.Sprintf(
		`INSERT INTO
			setting (user_id, imap_host, imap_port, smtp_host, smtp_port) 
		VALUES  
			(%d, '%s', %d, '%s', %d) 
		ON CONFLICT 
			(user_id) 
		DO UPDATE SET
			imap_host = '%s',
			imap_port = %d,
			smtp_host = '%s',
			smtp_port = %d`, 
		data.UserId, 
		data.IMAPHost, data.IMAPPort, data.SMTPHost, data.SMTPPort,
		data.IMAPHost, data.IMAPPort, data.SMTPHost, data.SMTPPort,
	)

	_, err := s.conn.Exec(query)
	if err != nil {
		log.Println(err.Error())
	}

	return err
}

func (s *SettingRepository) GetIMAPSettig(uid uint) (*dto.MailClient, error) {

	var mc dto.MailClient

	query := fmt.Sprintf(
		`SELECT 
			user_.username as username,
			user_.password as password,
			setting.imap_host as host,
			setting.imap_port as port
		FROM 
			user_
		LEFT JOIN
			setting ON setting.user_id = user_.id
		WHERE 
			user_.id = %d`, uid,
	)

	err := s.conn.Get(&mc, query)
	if err != nil {
		log.Println(err.Error())
	}

	return &mc, err
}

func (s *SettingRepository) GetSMTPSettig(uid uint) (*dto.MailClient, error) {

	var mc dto.MailClient

	query := fmt.Sprintf(
		`SELECT 
			user_.username as username,
			user_.password as password,
			setting.smtp_host as host,
			setting.smtp_port as port
		FROM 
			user_
		LEFT JOIN
			setting ON setting.user_id = user_.id
		WHERE 
			user_.id = %d`, uid,
	)

	err := s.conn.Get(&mc, query)
	if err != nil {
		log.Println(err.Error())
	}

	return &mc, err
}


