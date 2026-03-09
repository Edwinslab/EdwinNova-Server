package mail

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/go-resty/resty/v2"
)

type SMTPClient struct {
	host string
	port string
	auth smtp.Auth
	from string
}

func NewSMTPClient() (*SMTPClient, error) {

	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")

	if host == "" || port == "" || user == "" || pass == "" {
		return nil, fmt.Errorf("missing smtp environment variables")
	}

	auth := smtp.PlainAuth("", user, pass, host)

	return &SMTPClient{
		host: host,
		port: port,
		auth: auth,
		from: user,
	}, nil
}

func (s *SMTPClient) Send(to, subject, body string) error {

	msg := []byte(
		"Subject: " + subject + "\r\n" +
			"MIME-version: 1.0;\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
			body,
	)

	return smtp.SendMail(
		s.host+":"+s.port,
		s.auth,
		s.from,
		[]string{to},
		msg,
	)
}



func (s *SMTPClient) SendGenericEmail(toEmail, subject, body string) error {

	client := resty.New()

	payload := map[string]interface{}{
		"sender": map[string]string{
			"name":  "Linkly Media",
			"email": os.Getenv("BREVO_EMAIL"),
		},
		"to": []map[string]string{
			{"email": toEmail},
		},
		"subject": subject,
		"htmlContent": fmt.Sprintf(`
			<p>%s</p>
		`, body),
	}

	resp, err := client.R().
		SetHeader("accept", "application/json").
		SetHeader("api-key", os.Getenv("BREVO_API_KEY")).
		SetHeader("content-type", "application/json").
		SetBody(payload).
		Post("https://api.brevo.com/v3/smtp/email")

	if err != nil {
		return err
	}

	if resp.StatusCode() >= 300 {
		return fmt.Errorf("brevo error: %s", resp.String())
	}

	fmt.Println("Generic email sent via Brevo API")
	return nil
}