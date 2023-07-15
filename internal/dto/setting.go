package dto

type Setting struct {
	UserId		uint	`db:"user_id" json:"-"`
	SMTPHost	string	`json:"smtp_host" binding:"required" db:"smtp_host"`
	IMAPHost	string	`json:"imap_host" binding:"required" db:"imap_host"`
	SMTPPort	uint	`json:"smtp_port" binding:"required" db:"smtp_port"`
	IMAPPort	uint	`json:"imap_port" binding:"required" db:"imap_port"`
}

type MailClient struct {
	Username	string	`db:"username"`
	Password	string	`db:"password"`
	Host		string	`db:"host"`
	Port		uint	`db:"port"`
}
