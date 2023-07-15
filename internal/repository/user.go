package repository

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/amaretur/mail-client/internal/dto"
)

type UserRepository struct {
	conn *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		conn: db,
	}
}

func (u *UserRepository) GetById(id uint) (*dto.User, error) {

	var user dto.User

	query := fmt.Sprintf(
		"SELECT id, username, password FROM user_ WHERE id = %d",
		id,
	)

	err := u.conn.Get(&user, query)
	if err != nil {
		log.Println("")
		log.Println(err.Error())
	}

	return &user, err
}

func (u *UserRepository) CreateOrUpdate(username, password string) (uint, error) {

	query := fmt.Sprintf(
		`INSERT INTO
			user_ (username, password) 
		VALUES  
			('%s', '%s') 
		ON CONFLICT 
			(username) 
		DO UPDATE 
			SET password='%s'
		RETURNING id`, 
		username, password, password,
	)

	var uid uint

	err := u.conn.Get(&uid, query)
	if err != nil {
		log.Println(err.Error())
	}

	return uid, err
}

func (u *UserRepository) Delete(uid uint) error {

	_, err := u.conn.Exec(fmt.Sprintf("DELETE FROM user_ WHERE id = %d", uid))
	if err != nil {
		log.Println(err.Error())
	}

	return err
}
