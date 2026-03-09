package mail

import (
	"crypto/tls"
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
			"MIME-version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n" +
			body,
	)

	server := s.host + ":" + s.port

	tlsconfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         s.host,
	}

	conn, err := tls.Dial("tcp", server, tlsconfig)
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return err
	}

	if err = client.Auth(s.auth); err != nil {
		return err
	}

	if err = client.Mail(s.from); err != nil {
		return err
	}

	if err = client.Rcpt(to); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	client.Quit()

	return nil
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