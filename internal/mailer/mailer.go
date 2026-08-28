package mailer

import (
	"bytes"
	"crypto/tls"
	"embed"
	
	"io"
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
	// Use standard Go DNS resolution
	dialHost := host

	dialer := mail.NewDialer(dialHost, port, username, password)
	dialer.Timeout = 10 * time.Second
	dialer.LocalName = "teks-invoice.com" // Provide a valid EHLO name
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

func (m Mailer) SendMailWithAttachment(recipient, templateFile string, data interface{}, fileName string, fileData []byte) error {
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
	
	msg.Attach(fileName, mail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(fileData)
		return err
	}))

	err = m.dialer.DialAndSend(msg)
	if err != nil {
		return err
	}
	return nil
}
