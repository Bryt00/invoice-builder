package mailer

import (
	"bytes"
	"context"
	"crypto/tls"
	"embed"
	"net"
	"text/template"
	"time"

	"github.com/go-mail/mail"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

type Mailer struct {
	dialer *mail.Dialer
	sender string
}

func NewMailer(host string, port int, username, password, sender string) Mailer {
	dialHost := host
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{Timeout: 3 * time.Second}
			return d.DialContext(ctx, "udp", "8.8.8.8:53")
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if ips, err := resolver.LookupIP(ctx, "ip4", host); err == nil {
		for _, ip := range ips {
			if !ip.IsLoopback() {
				dialHost = ip.String()
				break
			}
		}
	}

	dialer := mail.NewDialer(dialHost, port, username, password)
	dialer.Timeout = 10 * time.Second
	dialer.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	return Mailer{
		dialer: dialer,
		sender: sender,
	}
}
func (m Mailer) SendMail(recipient, templateFile string, data interface{}) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+templateFile)

	if err != nil {
		return err
	}
	subject := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(subject, "subject", data)
	if err != nil {
		return err
	}
	body := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(body, "body", data)
	if err != nil {
		return err
	}
	htmlBody := new(bytes.Buffer)
	err = tmpl.ExecuteTemplate(htmlBody, "htmlBody", data)
	if err != nil {
		return err
	}
	msg := mail.NewMessage()
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", body.String())
	msg.SetBody("text/html", htmlBody.String())

	err = m.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}
	return nil
}
