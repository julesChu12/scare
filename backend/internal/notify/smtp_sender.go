package notify

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"time"
)

type SMTPSender struct {
	host     string
	port     int
	username string
	password string
	useTLS   bool
}

func NewSMTPSender(host string, port int, username, password string, useTLS bool) (*SMTPSender, error) {
	if host == "" || port == 0 {
		return nil, errors.New("smtp host/port required")
	}
	return &SMTPSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
		useTLS:   useTLS,
	}, nil
}

func (s *SMTPSender) Send(ctx context.Context, to, subject, body string) error {
	if to == "" {
		return errors.New("email recipient required")
	}
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	msg := buildEmailMessage(s.username, to, subject, body)
	if s.useTLS {
		return s.sendTLS(ctx, addr, msg, to)
	}

	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	return smtp.SendMail(addr, auth, s.username, []string{to}, []byte(msg))
}

func (s *SMTPSender) sendTLS(ctx context.Context, addr, msg, to string) error {
	dialer := &net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()

	deadline, ok := ctx.Deadline()
	if ok {
		_ = conn.SetDeadline(deadline)
	} else {
		_ = conn.SetDeadline(time.Now().Add(10 * time.Second))
	}

	tlsConn := tls.Client(conn, &tls.Config{ServerName: s.host})
	client, err := smtp.NewClient(tlsConn, s.host)
	if err != nil {
		return err
	}
	defer client.Quit()

	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	if err := client.Auth(auth); err != nil {
		return err
	}
	if err := client.Mail(s.username); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}
	writer, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := writer.Write([]byte(msg)); err != nil {
		_ = writer.Close()
		return err
	}
	return writer.Close()
}

func buildEmailMessage(from, to, subject, body string) string {
	headers := ""
	headers += fmt.Sprintf("From: %s\r\n", from)
	headers += fmt.Sprintf("To: %s\r\n", to)
	headers += fmt.Sprintf("Subject: %s\r\n", subject)
	headers += "MIME-Version: 1.0\r\n"
	headers += "Content-Type: text/plain; charset=UTF-8\r\n\r\n"
	return headers + body
}
