package mail

import (
	"errors"
	"fmt"
	"crypto/tls"
	"net/mail"
	"net/smtp"
)

type Mail struct {
	FromAddress string
	ToAddress   string
	Password    string
}

func (m Mail) Send(title, body string) error {
	from := mail.Address{"", m.FromAddress}
	to := mail.Address{"", m.ToAddress}
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", m.FromAddress, m.ToAddress, title, body)
	host, port := GetServerAddr(m.FromAddress)
	if host == "" {
		return errors.New("Not support this mail host")
	}
	server := host+":"+port
	auth := smtp.PlainAuth(
		"",
		from.String(),
		m.Password,
		host,
	)
	tlscfg := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}
	conn, err := tls.Dial("tcp", server, tlscfg)
	if err != nil {
		return err
	}
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer client.Quit()
	if err = client.Auth(auth); err != nil {
		return err
	}
	if err = client.Mail(from.Address); err != nil {
		return err
	}
	if err = client.Rcpt(to.Address); err != nil {
		return err
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	defer w.Close()
	if _, err = w.Write([]byte(message)); err != nil {
		return err
	}
	return nil
}