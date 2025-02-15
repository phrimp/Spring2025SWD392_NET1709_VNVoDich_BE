package services

import (
	"fmt"
	"google-service/internal/config"
	"net/smtp"
	"strings"
)

type EmailService struct {
	config *config.SMTPConfig
}

func NewEmailService(_config *config.SMTPConfig) *EmailService {
	return &EmailService{config: _config}
}

func (e *EmailService) SendEmail(title, body string, receipient []string) error {
	to := receipient

	message := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", strings.Join(to, ","), title, body))

	auth := smtp.PlainAuth("", e.config.Username, e.config.Password, e.config.Host)

	addr := fmt.Sprintf("%s:$s", e.config.Host, e.config.Port)

	err := smtp.SendMail(addr, auth, e.config.From, to, message)
	if err != nil {
		fmt.Printf("Error sending email: %s", err)
		return err
	}
	return nil
}
