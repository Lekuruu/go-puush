package email

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/smtp"

	"github.com/caarlos0/env/v11"
)

// SMTPConfig describes the configuration required for SMTP delivery
type SMTPConfig struct {
	From          string `env:"EMAIL_FROM"`
	Host          string `env:"SMTP_HOST"`
	Port          int    `env:"SMTP_PORT" envDefault:"587"`
	Username      string `env:"SMTP_USERNAME"`
	Password      string `env:"SMTP_PASSWORD"`
	EnableTLS     bool   `env:"SMTP_ENABLE_TLS" envDefault:"true"`
	SkipTLSVerify bool   `env:"SMTP_SKIP_TLS_VERIFY" envDefault:"false"`
}

func LoadSMTPConfig() (*SMTPConfig, error) {
	var config SMTPConfig
	if err := env.Parse(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

// SMTPEmail delivers messages using an SMTP server
type SMTPEmail struct {
	config *SMTPConfig
	auth   smtp.Auth
}

// NewSMTPEmail constructs an SMTP-backed email sender
func NewSMTPEmail(config *SMTPConfig) Email {
	return &SMTPEmail{config: config}
}

// FromAddress returns the configured default sender address.
func (s *SMTPEmail) FromAddress() string {
	return s.config.From
}

// Setup validates the SMTP configuration and prepares any required auth
func (s *SMTPEmail) Setup() error {
	if s.config.Host == "" {
		return errors.New("email: SMTP host is required")
	}

	if s.config.Port == 0 {
		s.config.Port = 587
	}

	if s.config.From == "" {
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

	mimeMessage, err := message.BuildMimeMessage(s.config.From)
	if err != nil {
		return err
	}

	address := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)
	client, err := smtp.Dial(address)
	if err != nil {
		return fmt.Errorf("email: failed to connect to SMTP server: %w", err)
	}
	defer client.Close()

	if s.config.EnableTLS {
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

	if err := client.Mail(s.config.From); err != nil {
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
