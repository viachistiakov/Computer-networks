package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
)

type Request struct {
	Command string      `json:"command"`
	Data    interface{} `:"data,omitempty"`
}

type Response struct {
	Status string           `json:"status"`
	Data   *json.RawMessage `json:"data,omitempty"`
}

func sendRequest(conn net.Conn, command string, data interface{}) error {
	request := Request{
		Command: command,
		Data:    data,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return err
	}

	_, err = conn.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

func receiveResponse(conn net.Conn) (*Response, error) {
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, err
	}

	var response Response
	err = json.Unmarshal(buffer[:n], &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func interact(conn net.Conn) {
	defer conn.Close()

	for {
		fmt.Print("command = ")
		var command string
		_, err := fmt.Scanln(&command)
		if err != nil {
			fmt.Println("error:", err)
			continue
		}

		switch command {
		case "quit":
			err := sendRequest(conn, "quit", nil)
			if err != nil {
				fmt.Println("error:", err)
			}
			return
		case "path":
			fmt.Print("path = ")
			var path string
			_, err := fmt.Scanln(&path)
			if err != nil {
				fmt.Println("error:", err)
				continue
			}

			err = sendRequest(conn, "path", path)
			if err != nil {
				fmt.Println("error:", err)
				continue
			}
		default:
			fmt.Println("error: unknown command")
			continue
		}

		response, err := receiveResponse(conn)
		if err != nil {
			fmt.Println("error:", err)
			break
		}

		switch response.Status {
		case "ok":
			fmt.Println("ok")
		case "failed":
			if response.Data == nil {
				fmt.Println("error: data field is absent in response")
			} else {
				var errorMsg string
				err := json.Unmarshal(*response.Data, &errorMsg)
				if err != nil {
					fmt.Println("error: malformed data field in response")
				} else {
					fmt.Println("failed:", errorMsg)
				}
			}
		case "result":
			if response.Data == nil {
				fmt.Println("error: data field is absent in response")
			} else {
				var files []string
				err := json.Unmarshal(*response.Data, &files)
				if err != nil {
					fmt.Println("error: malformed data field in response")
				} else {
					for _, file := range files {
						fmt.Println(file)
					}
				}
			}
		default:
			fmt.Printf("error: server reports unknown status %q\n", response.Status)
		}
	}
}

func main() {
	addr := flag.String("addr", "127.0.0.1:6000", "specify ip address and port")
	flag.Parse()

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	interact(conn)
}
