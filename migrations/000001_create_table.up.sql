BEGIN;

	DROP TABLE IF EXISTS user_;
	CREATE TABLE IF NOT EXISTS user_ (
		id SERIAL
			NOT NULL
			PRIMARY KEY,
		username VARCHAR(255)
			NOT NULL
			UNIQUE,
		password varchar(255)
			NOT NULL,
		verify_pubkey TEXT
			DEFAULT '',
		verify_privkey TEXT
			DEFAULT ''
	);

	DROP TABLE IF EXISTS setting;
	CREATE TABLE IF NOT EXISTS setting (
		id SERIAL
			NOT NULL
			PRIMARY KEY,
		user_id INTEGER
			UNIQUE
			REFERENCES user_(id)
			ON UPDATE CASCADE
			ON DELETE CASCADE,
		imap_host VARCHAR(255)
			NOT NULL,
		imap_port INTEGER
			NOT NULL,
		smtp_host VARCHAR(255)
			NOT NULL,
		smtp_port INTEGER
			NOT NULL
	);

	DROP TABLE IF EXISTS dialog;
	CREATE TABLE IF NOT EXISTS dialog (
		id SERIAL
			NOT NULL
			PRIMARY KEY,
		user_id INTEGER
			REFERENCES user_(id)
				ON UPDATE CASCADE
				ON DELETE CASCADE,
		interlocutor VARCHAR(255)
			NOT NULL,
		encrypt_key TEXT
			DEFAULT '',
		decrypt_key TEXT
			DEFAULT '',
		verify_key TEXT
			DEFAULT '',
		share_encrypt_key TEXT
			DEFAULT '',

		UNIQUE(user_id, interlocutor)
	);

COMMIT;


