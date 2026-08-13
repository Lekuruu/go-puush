package config

type EmailConfig struct {
	Type string `env:"EMAIL_TYPE" envDefault:"noop"`
	From string `env:"EMAIL_FROM"`
	SMTP SMTPConfig
}

// SMTPConfig describes the configuration required for SMTP delivery
type SMTPConfig struct {
	Host          string `env:"SMTP_HOST"`
	Port          int    `env:"SMTP_PORT" envDefault:"587"`
	Username      string `env:"SMTP_USERNAME"`
	Password      string `env:"SMTP_PASSWORD"`
	UseTLS        bool   `env:"SMTP_USE_TLS" envDefault:"true"`
	SkipTLSVerify bool   `env:"SMTP_SKIP_TLS_VERIFY" envDefault:"false"`
}
