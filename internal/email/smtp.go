package email

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/smtp"

	"github.com/Lekuruu/go-puush/internal/config"
)

// SMTPEmail delivers messages using an SMTP server
type SMTPEmail struct {
	from   string
	config config.SMTPConfig
	auth   smtp.Auth
}

// NewSMTPEmail constructs an SMTP-backed email sender
func NewSMTPEmail(from string, smtpConfig config.SMTPConfig) Email {
	return &SMTPEmail{from: from, config: smtpConfig}
}

// FromAddress returns the configured default sender address.
func (s *SMTPEmail) FromAddress() string {
	return s.from
}

// Setup validates the SMTP configuration and prepares any required auth
func (s *SMTPEmail) Setup() error {
	if s.config.Host == "" {
		return errors.New("email: SMTP host is required")
	}

	if s.config.Port == 0 {
		s.config.Port = 587
	}

	if s.from == "" {
		return errors.New("email: SMTP from address is required")
	}

	if s.config.Username != "" {
		s.auth = smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	}

	return nil
}

// Send delivers the provided message using SMTP
func (s *SMTPEmail) Send(message *Message) error {
	if err := message.Validate(); err != nil {
		return err
	}

	mimeMessage, err := message.BuildMimeMessage(s.from)
	if err != nil {
		return err
	}

	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	client, err := smtp.Dial(address)
	if err != nil {
		return fmt.Errorf("email: failed to connect to SMTP server: %w", err)
	}
	defer client.Close()

	if s.config.UseTLS {
		if ok, _ := client.Extension("STARTTLS"); ok {
			tlsConfig := &tls.Config{
				ServerName:         s.config.Host,
				InsecureSkipVerify: s.config.SkipTLSVerify,
			}
			if err := client.StartTLS(tlsConfig); err != nil {
				return fmt.Errorf("email: failed to start TLS: %w", err)
			}
		}
	}

	if s.auth != nil {
		if ok, _ := client.Extension("AUTH"); ok {
			if err := client.Auth(s.auth); err != nil {
				return fmt.Errorf("email: failed to authenticate: %w", err)
			}
		}
	}

	if err := client.Mail(s.from); err != nil {
		return fmt.Errorf("email: failed to set sender: %w", err)
	}

	for _, recipient := range message.To {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("email: failed to add recipient %s: %w", recipient, err)
		}
	}

	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("email: failed to open SMTP data writer: %w", err)
	}

	if _, err := writer.Write(mimeMessage); err != nil {
		writer.Close()
		return fmt.Errorf("email: failed to write message body: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("email: failed to finalize message body: %w", err)
	}

	if err := client.Quit(); err != nil {
		return fmt.Errorf("email: failed to close SMTP connection: %w", err)
	}

	return nil
}
