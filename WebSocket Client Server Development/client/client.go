package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	utils "lab5"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8000", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		for {
			scanner := bufio.NewScanner(os.Stdin)
			var seq utils.Sequence
			fmt.Println("Введите количество чисел, которые хотите ввести: ")
			var n int
			scanner.Scan()
			n, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Println("Invalid number:", err)
				return
			}
			fmt.Println("Введите Ваши действительные числа(каждое на новой строчке): ")
			for i := 0; i < n; i++ {
				scanner.Scan()
				num, err := strconv.ParseFloat(strings.TrimSpace(scanner.Text()), 64)
				if err != nil {
					log.Println("Invalid number:", err)
					return
				}
				seq.Data = append(seq.Data, num)
			}

			message, err := json.Marshal(seq)
			if err != nil {
				log.Fatalln(err)
			}

			err = c.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("write:", err)
				return
			}
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("%s", message)
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
