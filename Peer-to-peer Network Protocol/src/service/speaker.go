package service

import (
	"encoding/json"
	"fmt"
	"net"
)

func (p *Peer) dial(msg Message) {
	conn, err := net.Dial("tcp", p.NextAddress())
	if err != nil {
		p.Logger.Println(err.Error())
		return
	}
	defer conn.Close()

	jBytes, err := json.Marshal(msg)
	if err != nil {
		p.Logger.Panic(err.Error())
	}

	_, err = conn.Write(jBytes)
	if err != nil {
		p.Logger.Panic(err.Error())
	}
	p.Logger.Printf("%s: sent message to %s successfully\n",
		p.Address(), p.NextAddress())
}

func (p *Peer) say() {
	for {
		var cmd string
		fmt.Print("\nEnter one of the commands (add/list/stop): ")
		fmt.Scan(&cmd)

		switch cmd {
		case "add":
			var seg Segment
			fmt.Printf("xA = ")
			fmt.Scan(&seg.A.X)
			fmt.Printf("yA = ")
			fmt.Scan(&seg.A.Y)
			fmt.Printf("xB = ")
			fmt.Scan(&seg.B.X)
			fmt.Printf("yB = ")
			fmt.Scan(&seg.B.Y)
			p.mu.Lock() // для поддержки асинхронности
			p.Segments = append(p.Segments, seg)
			p.mu.Unlock()
			fmt.Println("Success!")
			go p.dial(Message{p.Name, seg})
		case "list":
			if len(p.Segments) == 0 {
				fmt.Println("No lines!")
			}
			for _, seg := range p.Segments {
				fmt.Printf("line A={%v;%v}, B={%v;%v}\n",
					seg.A.X, seg.A.Y, seg.B.X, seg.B.Y)
			}
		case "stop":
			p.Stop <- struct{}{}
			return
		default:
			fmt.Println("Incorrect command! Try again.")
		}
	}
}
