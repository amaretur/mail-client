package repository

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	"github.com/amaretur/mail-client/internal/dto"
)

type KeyRepository struct {
	conn *sqlx.DB
}

func NewKeyRepository(db *sqlx.DB) *KeyRepository {
	return &KeyRepository{
		conn: db,
	}
}

func (k *KeyRepository) UpdateVerify(uid uint, keys *dto.VerifyKeys) error {

	query := fmt.Sprintf(
		`UPDATE user_ SET
			verify_pubkey = '%s', 
			verify_privkey = '%s'
		WHERE
			id = %d`, 
		keys.Public, keys.Private, uid,
	)

	_, err := k.conn.Exec(query)

	return err
}

func (k *KeyRepository) GetVerify(uid uint) (*dto.VerifyKeys, error) {

	query := fmt.Sprintf(
		"SELECT verify_pubkey, verify_privkey FROM user_ WHERE id = %d", uid)

	var data dto.VerifyKeys

	err := k.conn.Get(&data, query)
	if err != nil {
		log.Println(err.Error())
	}

	return &data, err
}

func (k *KeyRepository) GetDialogKeysList(
	uid, offset, limit uint,
) (*dto.DialogKeysList, error) {

	query := fmt.Sprintf(
		`SELECT
			id, interlocutor 
		FROM 
			dialog 
		WHERE 
			user_id = %d
		ORDER BY 
			interlocutor
		OFFSET %d LIMIT %d`, 
		uid, offset, limit,
	)

	var keys dto.DialogKeysList

	err := k.conn.Select(&keys.Sets, query)

	return &keys, err
}

func (k *KeyRepository) GetDialogKeys(id uint) (*dto.DialogKeys, error) {

	query := fmt.Sprintf(
		`SELECT
			encrypt_key,
			decrypt_key,
			verify_key,
			share_encrypt_key
		FROM 
			dialog
		WHERE 
			id = %d`,
		id,
	)

	var keys dto.DialogKeys

	err := k.conn.Get(&keys, query)

	return &keys, err
} 

func (k *KeyRepository) CreateOrUpdateInterlocutorKeys(
	uid uint, interlocutor, encrypt, verify string,
) (uint, error) {

	query := fmt.Sprintf(`
		INSERT INTO
			dialog (interlocutor, user_id, encrypt_key, verify_key)
		VALUES
			('%s', %d, '%s', '%s')
		ON CONFLICT 
			(user_id, interlocutor)
		DO UPDATE SET
			encrypt_key = '%s',
			verify_key = '%s'
		RETURNING 
			id`,

		interlocutor, uid, encrypt, verify,
		encrypt, verify,
	)

	var id uint

	err := k.conn.Get(&id, query)
	if err != nil {
		log.Println(err)
	}

	return id, err
}

func (k *KeyRepository) CreateOrUpdateUserKeys(
	uid uint, interlocutor, decrypt, encrypt string,
) (uint, error) {

	query := fmt.Sprintf(`
		INSERT INTO
			dialog (interlocutor, user_id, decrypt_key, share_encrypt_key)
		VALUES
			('%s', %d, '%s', '%s')
		ON CONFLICT 
			(user_id, interlocutor)
		DO UPDATE SET
			decrypt_key = '%s',
			share_encrypt_key = '%s'
		RETURNING 
			id`,

		interlocutor, uid, decrypt, encrypt,
		decrypt, encrypt,
	)

	var id uint

	err := k.conn.Get(&id, query)

	return id, err
}

func (k *KeyRepository) GetDecryptAndVerifyKeys(
	uid uint, interlocutor string) (string, string, error) {

	query := fmt.Sprintf(`
		SELECT
			decrypt_key,
			verify_key
		FROM
			dialog
		WHERE
			user_id = %d AND interlocutor = '%s'
		`,

		uid, interlocutor,
	)

	var keys struct {
		DecryptKey	string	`db:"decrypt_key"`
		VerifyKey	string	`db:"verify_key"`
	}

	err := k.conn.Get(&keys, query)

	return keys.DecryptKey, keys.VerifyKey, err
}

func (k *KeyRepository) GetEncryptAndSignKeys(
	uid uint, interlocutor string) (string, string, error) {

	query := fmt.Sprintf(`
		SELECT
			dialog.encrypt_key as encrypt,
			user_.verify_privkey as sign
		FROM
			dialog
		LEFT JOIN
			user_ ON user_.id = dialog.user_id
		WHERE
			user_id = %d AND interlocutor = '%s'`,

		uid, interlocutor,
	)

	var keys struct {
		EncryptKey	string	`db:"encrypt"`
		SignKey	string		`db:"sign"`
	}

	err := k.conn.Get(&keys, query)

	return keys.EncryptKey, keys.SignKey, err
}
