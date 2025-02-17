package config

type EmailConfig struct {
	SMTPConfig SMTPConfig
}

func NewEmailConfig() *EmailConfig {
	return &EmailConfig{
		SMTPConfig: loadSMTPConfig(),
	}
}
