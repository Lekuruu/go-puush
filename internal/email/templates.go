package email

const WelcomeTemplate = `Welcome to puush!

You are now ready to login to your puush app(s) using your new account. Your username is your email address.

Sincerely,
The puush Team

(If you didn't request this, simply ignore this email and we won't contact you again)`

func FormatWelcomeEmail(to string) *Message {
	return &Message{
		To:       []string{to},
		Subject:  "[puush] Account Details",
		TextBody: WelcomeTemplate,
	}
}
