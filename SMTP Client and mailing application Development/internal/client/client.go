package client

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"lab6_2/internal/model"
	"log"
	"net/smtp"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	withoutError = false
	withError    = true
)

type Repository interface {
	GetUsers() ([]*model.User, error)
	SaveLog(string, string, bool) error
}

type Client struct {
	repo     Repository
	from     string
	pass     string
	hostSMTP string
	portSMTP int
}

func NewClient(repo Repository, from, password, hostSMTP string, portSMTP int) *Client {
	return &Client{
		repo:     repo,
		from:     from,
		pass:     password,
		hostSMTP: hostSMTP,
		portSMTP: portSMTP,
	}
}

func (c *Client) SendMailList() {
	log.Println("App started, getting users list from repository...")
	users, err := c.repo.GetUsers()
	if err != nil {
		log.Printf("Error occurred while getting users list: %v\n", err)
		return
	}
	log.Println("Got users list from repository, mailing started...")
	for _, u := range users {
		err := c.sendEmail(u)
		if err != nil {
			log.Println(err)
			err = c.repo.SaveLog(u.Email, fmt.Sprintf("Error occurred while "+
				"sending email: %v", err), withError)
			if err != nil {
				log.Println(err)
			}
		} else {
			log.Printf("Email to %v sent successfully\n", u.Email)
			err = c.repo.SaveLog(u.Email, fmt.Sprintf("Successfully sent email"),
				withoutError)
			if err != nil {
				log.Println(err)
			}
		}
	}
	log.Printf("Sent all emails at %s!\n", time.Now().Format("2006-01-02 15:04:05"))
}

func getTemplate() (*template.Template, error) {
	tpl, err := template.ParseFiles("internal/templates/email.html")
	if err != nil {
		return nil, err
	}
	return tpl, nil
}

func (c *Client) sendEmail(user *model.User) error {
	tlsConfig := &tls.Config{InsecureSkipVerify: true}

	serverAddr := fmt.Sprintf("%s:%d", c.hostSMTP, c.portSMTP)

	auth := smtp.PlainAuth("", c.from, c.pass, c.hostSMTP)
	conn, err := tls.Dial("tcp", serverAddr, tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Send EHLO to establish connection
	client, err := smtp.NewClient(conn, c.hostSMTP)
	if err != nil {
		return err
	}
	defer client.Quit()

	if err := client.Auth(auth); err != nil {
		return err
	}

	if err := client.Mail(c.from); err != nil {
		return err
	}
	if err := client.Rcpt(user.Email); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	defer w.Close()

	headers := fmt.Sprintf("Subject: Greetings!\r\n"+
		"To: %s\r\n"+
		"Content-Type: text/html\r\n", user.Email)

	if _, err := w.Write([]byte(headers)); err != nil {
		return err
	}

	tpl, err := getTemplate()
	if err != nil {
		return err
	}

	if err := tpl.Execute(w, user); err != nil {
		return err
	}

	return nil
}
