package main

import (
	"lab6_1/client"
	"lab6_1/client/terminal"
	"log"
	"time"

	"github.com/jlaffaye/ftp"
)

func main() {
	clientFTP, err := ftp.Dial("185.139.70.64:55321", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	err = clientFTP.Login("user", "1234test")
	if err != nil {
		log.Fatal(err)
	}

	c, err := client.NewClient(clientFTP, "./downloads")
	if err != nil {
		log.Fatal(err)he
	}
	t := terminal.NewTerminal(c)
	t.Listen()

	if err := clientFTP.Quit(); err != nil {
		log.Fatal(err)
	}

}
