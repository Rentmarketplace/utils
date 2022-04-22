package utils

import (
	"github.com/xhit/go-simple-mail/v2"
	"os"
	"strconv"
	"time"
)

type MailServer struct{}

type Connection struct {
	Port       int
	Encryption int
	Host       string
	Username   string
	Password   string
}

func (m *MailServer) newConnection() (*Connection, error) {
	port, err := strconv.Atoi(os.Getenv("EMAIL_PORT"))

	if err != nil {
		Logger().Errorln(err)
		return nil, err
	}

	encrypt, err := strconv.Atoi(os.Getenv("EMAIL_ENCRYPTION"))

	if err != nil {
		Logger().Errorln(err)
		return nil, err
	}

	return &Connection{
		Port:       port,
		Encryption: encrypt,
		Host:       os.Getenv("EMAIl_HOST"),
		Username:   os.Getenv("EMAIL_USERNAME"),
		Password:   os.Getenv("EMAIL_PASSWORD"),
	}, nil
}

func (m *MailServer) EmailSettings() (*mail.SMTPClient, error) {
	server := mail.NewSMTPClient()

	con, err := m.newConnection()

	if err != nil {
		return nil, err
	}

	server.Port = con.Port
	server.Host = con.Host
	server.Username = con.Username
	server.Password = con.Password
	server.Encryption = mail.Encryption(con.Encryption)

	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second

	smtpClient, _ := server.Connect()
	return smtpClient, nil
}
