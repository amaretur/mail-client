package dto

type DialogKeys struct {
	Id				uint	`json:"id" db:"id"`
	Verify			string	`json:"verify" db:"verify_key"`
	Encrypt			string	`json:"encrypt" db:"encrypt_key"`
	Decrypt			string	`json:"decrypt" db:"decrypt_key"`
	ShareEncrypt	string	`json:"shared" db:"share_encrypt_key"`
}

type DialogKeysBrief struct {
	Id				uint	`json:"id" db:"id"`
	Interlocutor	string	`json:"interlocutor" db:"interlocutor"`
}

type DialogKeysList struct {
	Sets []DialogKeysBrief	`json:"sets"`
}

type VerifyKeys struct {
	Public	string	`json:"public" db:"verify_pubkey"`
	Private	string	`json:"private" db:"verify_privkey"`
}

type EncKeys struct {
	Public	string	`json:"public" db:"encrypt_key"`
	Private	string	`json:"private" db:"decrypt_key"`
}

type Share struct {
	Encrypt	string	`json:"encrypt"`
	Verify	string	`json:"verify"`
}

type Receive struct {
	Interlocutor	string	`json:"interlocutor" binding:"required"`
	Encrypt			string	`json:"encrypt" binding:"required"`
	Verify			string	`json:"verify" binding:"required"`
}
