package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
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
			fmt.Println(err.Error())
			return
		}

		s.Stdout = os.Stdout
		s.Stderr = os.Stderr
		err = s.Run(cmd)
		if err != nil {
			fmt.Println("command execution failed")
			fmt.Println(err.Error())
		}
		s.Close()
	}
}

func main() {
	conf := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			ssh.Password("j78Ei7372PRf"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	var serverIP, serverPort string
	fmt.Print("Enter server IP: ")
	fmt.Scan(&serverIP)
	fmt.Print("Enter server port: ")
	fmt.Scan(&serverPort)

	conn, err := ssh.Dial("tcp", serverIP+":"+serverPort, conf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()

	c := &client{
		config:     conf,
		connection: conn,
	}
	c.start()
}
