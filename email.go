package pills

import (
	"fmt"
	"net/smtp"
	"strconv"

	Z "github.com/rwxrob/bonzai/z"
)

type SmtpConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
	To       string
}

func GetSmtpConfig(v *Z.Cmd) (*SmtpConfig, error) {
	var err error
	config := &SmtpConfig{}

	config.Host, err = v.Get("smtp.host")
	if err != nil {
		return nil, err
	}
	if config.Host == "" {
		return nil, fmt.Errorf("smtp.host is not set")
	}

	portStr, err := v.Get("smtp.port")
	if err != nil {
		return nil, err
	}
	if portStr == "" {
		return nil, fmt.Errorf("smtp.port is not set")
	}
	config.Port, err = strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid smtp.port: %w", err)
	}

	config.User, err = v.Get("smtp.user")
	if err != nil {
		return nil, err
	}
	if config.User == "" {
		return nil, fmt.Errorf("smtp.user is not set")
	}

	config.Password, err = v.Get("smtp.pass")
	if err != nil {
		return nil, err
	}
	if config.Password == "" {
		return nil, fmt.Errorf("smtp.pass is not set")
	}

	config.From, err = v.Get("smtp.from")
	if err != nil {
		return nil, err
	}
	if config.From == "" {
		return nil, fmt.Errorf("smtp.from is not set")
	}

	config.To, err = v.Get("smtp.to")
	if err != nil {
		return nil, err
	}
	if config.To == "" {
		return nil, fmt.Errorf("smtp.to is not set")
	}

	return config, nil
}

func (config *SmtpConfig) SendEmail(subject, body string) error {
	auth := smtp.PlainAuth("", config.User, config.Password, config.Host)
	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	msg := []byte("To: " + config.To + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")

	return smtp.SendMail(addr, auth, config.From, []string{config.To}, msg)
}
