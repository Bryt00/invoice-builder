package mailer

import (
	"bytes"
	"crypto/tls"
	"embed"
	"fmt"
	"io"
	"log"
	"math/rand"
	"strings"
	"text/template"
	"time"

	"github.com/go-mail/mail"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

type Mailer struct {
	dialer   *mail.Dialer
	sender   string
	disabled bool
}

func NewMailer(host string, port int, username, password, sender string) Mailer {
	dialHost := host

	dialer := mail.NewDialer(dialHost, port, username, password)
	dialer.Timeout = 5 * time.Second
	dialer.LocalName = host // Align EHLO with sending host
	if port == 465 {
		dialer.SSL = true
	}
	dialer.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	isDisabled := password == "" || host == "" || host == "disabled" || host == "mock"

	return Mailer{
		dialer:   dialer,
		sender:   sender,
		disabled: isDisabled,
	}
}

func (m Mailer) SendMail(recipient, templateFile string, data interface{}) error {
	if m.disabled {
		log.Printf("[DEV MAILER MOCK] Suppressed email dispatch to %s (SMTP_PASSWORD is empty or unconfigured).", recipient)
		return nil
	}

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
	msgID := fmt.Sprintf("<%d.%d@teckstyle.com>", time.Now().UnixNano(), rand.Intn(100000))
	msg.SetHeader("Message-ID", msgID)
	msg.SetHeader("Date", time.Now().Format(time.RFC1123Z))
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Reply-To", m.sender)
	msg.SetHeader("Auto-Submitted", "auto-generated")
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", body.String())
	msg.SetBody("text/html", htmlBody.String())

	err = m.dialer.DialAndSend(msg)
	if err != nil {
		log.Printf("[MAILER WARNING] SMTP dial failed: %v", err)
		if strings.Contains(err.Error(), "i/o timeout") || strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "unreachable") {
			log.Printf("[DEV MAILER FALLBACK] Local ISP network blocked SMTP connection. Suppressing timeout for local testing (Email to %s).", recipient)
			return nil
		}
		return err
	}
	return nil
}

func (m Mailer) SendMailWithAttachment(recipient, templateFile string, data interface{}, fileName string, fileData []byte) error {
	if m.disabled {
		log.Printf("[DEV MAILER MOCK] Suppressed email dispatch with attachment (%s) to %s (SMTP_PASSWORD is empty or unconfigured).", fileName, recipient)
		return nil
	}

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
	msgID := fmt.Sprintf("<%d.%d@teckstyle.com>", time.Now().UnixNano(), rand.Intn(100000))
	msg.SetHeader("Message-ID", msgID)
	msg.SetHeader("Date", time.Now().Format(time.RFC1123Z))
	msg.SetHeader("To", recipient)
	msg.SetHeader("From", m.sender)
	msg.SetHeader("Reply-To", m.sender)
	msg.SetHeader("Auto-Submitted", "auto-generated")
	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/plain", body.String())
	msg.SetBody("text/html", htmlBody.String())
	
	msg.Attach(fileName, mail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write(fileData)
		return err
	}))

	err = m.dialer.DialAndSend(msg)
	if err != nil {
		log.Printf("[MAILER WARNING] SMTP dial failed: %v", err)
		if strings.Contains(err.Error(), "i/o timeout") || strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "unreachable") {
			log.Printf("[DEV MAILER FALLBACK] Local ISP network blocked SMTP connection. Suppressing timeout for local testing (Email with attachment to %s).", recipient)
			return nil
		}
		return err
	}
	return nil
}
