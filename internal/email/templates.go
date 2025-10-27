package email

import "fmt"

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

// NOTE: This is not the original activation email, I still have to find one
const ActivationTemplate = `Welcome to puush!

Before you can start using your account, please confirm your email address by clicking the link below:

https://puush.me/account/activate?key=%s

Sincerely,
The puush Team

(If you didn't request this account, you can safely ignore this email.)`

func FormatActivationEmail(to string, activationKey string) *Message {
	return &Message{
		To:       []string{to},
		Subject:  "[puush] Activate Your Account",
		TextBody: fmt.Sprintf(ActivationTemplate, activationKey),
	}
}
