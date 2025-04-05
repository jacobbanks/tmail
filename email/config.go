package email

type Config struct {
	SMTPHost string
	SMTPPort string
	IMAPHost string
	IMAPPort string
}

var DefaultConfig = Config{
	SMTPHost: "smtp.gmail.com",
	SMTPPort: "587",
	IMAPHost: "imap.gmail.com",
	IMAPPort: "993",
}

func (c *Config) GetSMTPAddress() string {
	return c.SMTPHost + ":" + c.SMTPPort
}

func (c *Config) GetIMAPAddress() string {
	return c.IMAPHost + ":" + c.IMAPPort
}
