package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

const (
	IP       = "151.248.113.144"
	PORT     = "443"
	LOGIN    = "test"
	PASSWORD = "SDHBCXdsedfs222"
)

type client struct {
	config     *ssh.ClientConfig
	connection *ssh.Client
}

func (c *client) start() {
	fmt.Printf("Connected to %s\n", c.connection.RemoteAddr())
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		cmd := scanner.Text()

		if strings.Compare(cmd, "exit") == 0 {
			return
		}

		s, err := c.connection.NewSession()
		if err != nil {
			return
		}

		s.Stdout = os.Stdout
		s.Stderr = os.Stderr
		err = s.Run(cmd)
		if err != nil {
			fmt.Println("command execution failed")
		}
		s.Close()
	}
}

func main() {

	conf := &ssh.ClientConfig{
		User: LOGIN,
		Auth: []ssh.AuthMethod{
			ssh.Password(PASSWORD),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", IP+":"+PORT, conf)
	if err != nil {
		return
	}
	defer conn.Close()

	c := &client{
		config:     conf,
		connection: conn,
	}
	c.start()
}
