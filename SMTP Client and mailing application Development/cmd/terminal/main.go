package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

func main() {
	fmt.Print("To: ")
	to := readInput()

	fmt.Print("Subject: ")
	subject := readInput()

	fmt.Print("Message body: ")
	messageBody := readInput()

	username := "dts21@dactyl.su"
	password := "12345678990DactylSUDTS"

	smtpServer := "mail.nic.ru"
	smtpPort := 465

	message := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, messageBody)

	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         smtpServer,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", smtpServer, smtpPort), tlsConfig)
	if err != nil {
		fmt.Println("Error while connecting to SMTP:", err)
		return
	}
	defer conn.Close()

	auth := smtp.PlainAuth("", username, password, smtpServer)

	client, err := smtp.NewClient(conn, smtpServer)
	if err != nil {
		fmt.Println("Error while creating client session:", err)
		return
	}
	defer client.Quit()

	if err := client.Auth(auth); err != nil {
		fmt.Println("Auth error:", err)
		return
	}

	if err := client.Mail(username); err != nil {
		fmt.Println("Sender error:", err)
		return
	}

	if err := client.Rcpt(to); err != nil {
		fmt.Println("Recipient error:", err)
		return
	}

	writer, err := client.Data()
	if err != nil {
		fmt.Println("Data send error:", err)
		return
	}
	defer writer.Close()

	_, err = writer.Write([]byte(message))
	if err != nil {
		fmt.Println("Write data error:", err)
		return
	}

	fmt.Println("Message sent successfully.")
}

func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
