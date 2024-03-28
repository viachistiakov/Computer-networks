package terminal

import (
	"fmt"
	"lab6_1/client"
)

type Terminal struct {
	cl *client.Client
}

func NewTerminal(client *client.Client) *Terminal {
	return &Terminal{cl: client}
}

func (t *Terminal) Listen() {
	for {
		fmt.Print("Enter one of the commands (print \"help\" for help): ")
		var cmd string
		fmt.Scan(&cmd)
		switch cmd {
		case "list":
			fmt.Print("Enter directory path: ")
			var dir string
			fmt.Scan(&dir)
			s, err := t.cl.ListFiles(dir)
			if err != nil {
				fmt.Printf("Error occurred: %v\n", err)
			} else {
				fmt.Println("List of files:")
				fmt.Println(s)
			}
		case "load":
			fmt.Print("Enter filename: ")
			var name string
			fmt.Scan(&name)
			err := t.cl.LoadFile(name)
			if err != nil {
				fmt.Printf("Error occurred: %v\n", err)
			} else {
				fmt.Println("File was loaded successfully!")
			}
		case "send":
			fmt.Print("Enter filename: ")
			var name string
			fmt.Scan(&name)
			fmt.Print("Enter filename for saving file on server: ")
			var path string
			fmt.Scan(&path)
			err := t.cl.SendFile(name, path)
			if err != nil {
				fmt.Printf("Error occurred: %v\n", err)
			} else {
				fmt.Println("File was sent successfully!")
			}
		case "mkdir":
			fmt.Print("Enter directory name: ")
			var dir string
			fmt.Scan(&dir)
			err := t.cl.Mkdir(dir)
			if err != nil {
				fmt.Printf("Error occurred: %v\n", err)
			} else {
				fmt.Println("Directory was created successfully!")
			}
		case "del":
			fmt.Print("Enter filename: ")
			var name string
			fmt.Scan(&name)
			err := t.cl.DeleteFile(name)
			if err != nil {
				fmt.Printf("Error occurred: %v\n", err)
			} else {
				fmt.Println("File was deleted successfully!")
			}
		case "stop":
			fmt.Println("Shutting down...")
			return
		default:
			fmt.Println("Available commands:")
			fmt.Println("list  - list files in requested directory")
			fmt.Println("load  - load file from ftp server")
			fmt.Println("send  - send file to ftp server (with provided path)")
			fmt.Println("mkdir - make a directory on server")
			fmt.Println("del   - delete a file on server")
			fmt.Println("stop  - stop the client")
		}
	}
}
