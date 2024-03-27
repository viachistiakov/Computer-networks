package main

import (
	"fmt"
	"io"
	"lab4/server/service"
	"os/exec"
	"strings"

	"github.com/gliderlabs/ssh"
	"github.com/pkg/sftp"
	"golang.org/x/term"
)

// handleSFTP handles SFTP connection
func handleSFTP(sess ssh.Session) {
	server, err := sftp.NewServer(sess)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	if err := server.Serve(); err == io.EOF {
		server.Close()
		fmt.Println("sftp client exited session.")
	} else if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// handleSSH is the main handler for ssh
func handleSSH(s ssh.Session) {
	t := term.NewTerminal(s, "> ")
	for {
		line, err := t.ReadLine()
		if err == io.EOF {
			fmt.Println("ssh client exited session.")
			return
		} else if err != nil {
			fmt.Println(err.Error())
			continue
		}
		fields := strings.Fields(line)
		cmd := exec.Command(fields[0], fields[1:]...)
		cmd.Stdout = s
		cmd.Stderr = s
		if err := cmd.Run(); err != nil {
			fmt.Println(err.Error())
		}
	}
}

func main() {
	a := service.Auth{
		Users: map[string]string{
			"budnikov": "test1234",
		},
	}

	var port string
	fmt.Print("Enter port: ")
	fmt.Scan(&port)

	server := ssh.Server{
		Addr:            "0.0.0.0:" + port,
		Handler:         handleSSH,
		PasswordHandler: a.HandlePassword,
		SubsystemHandlers: map[string]ssh.SubsystemHandler{
			"sftp": handleSFTP,
		},
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err.Error())
		return
	}
}
