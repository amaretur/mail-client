package mails

type Mailbox struct {
	Id			string	`json:"id"`
	Prefix		string	`json:"prefix"`
	Delimiter	string	`json:"delimiter"`
	Name		string	`json:"name"`
}

func (m *Mailbox) String() string {
	return m.Prefix + m.Delimiter + m.Name
}

type File struct {
	Name	string	`json:"name"`
	IsHTML	bool	`json:"is_html"`
	Type	string	`json:"type"`
	Data	string	`json:"data"`
}

type Mail struct {
	Id			uint32	`json:"id"`
	Subject		string	`json:"subject"`
	FromName	string	`json:"from_name"`
	FromMail	string	`json:"from_mail"`
	ToName		string	`json:"to_name"`
	ToMail		string	`json:"to_mail"`
	Date		string	`json:"date"`
	Time		string	`json:"time"`
	Files		[]File	`json:"files"`
	Error		string	`json:"error,omitempty"`
}

func (m *Mail) GetHTML() string {
	for _, file := range m.Files {
		if file.IsHTML {
			return file.Data
		}
	}

	return ""
}

func (m *Mail) GetFiles() []File {
	res := []File{}

	for _, file := range m.Files {
		if !file.IsHTML {
			res = append(res, file)
		}
	}

	return res
}

type Protocol struct {
	Server	string	`json:"server"`
	Port	int		`json:"port"`
}

type Client struct {
	Login		string		`json:"login"`
	Password	string		`json:"password"`
	SMTP		Protocol	`json:"smtp"`
	IMAP		Protocol	`json:"imap"`
}


