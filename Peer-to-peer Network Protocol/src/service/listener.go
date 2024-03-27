package service

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
)

func (p *Peer) handleConnection(c net.Conn) {
	defer c.Close()
	d := json.NewDecoder(c)

	var msg Message
	err := d.Decode(&msg)
	if err != nil {
		p.Logger.Panic(err.Error())
	}

	p.Logger.Printf("%s: got message from %s\n",
		p.Address(), msg.SourceName)
	a := msg.Segment.A
	b := msg.Segment.B
	fmt.Print("Line length: ")
	fmt.Println(math.Sqrt(float64(float64(b.X-a.X)*float64(b.X-a.X) + float64(b.Y-a.Y)*float64(b.Y-a.Y))))
	if msg.SourceName == p.Name {
		return
	}
	go p.dial(Message{msg.SourceName, msg.Segment})
	p.mu.Lock() // для поддержки асинхронности
	p.Segments = append(p.Segments, msg.Segment)
	p.mu.Unlock()

	fmt.Printf("%s: read message from %s\n",
		p.Address(), msg.SourceName)
	p.Logger.Printf("%s: read message from %st\n",
		p.Address(), msg.SourceName)
}

func (p *Peer) listen() {
	l, err := net.Listen("tcp", p.Address())
	if err != nil {
		p.Logger.Println(err.Error())
		p.Stop <- struct{}{}
		panic(err.Error())
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			select {
			case <-p.Stop:
				return
			default:
				p.Logger.Panic(err.Error())
			}
		} else {
			go p.handleConnection(conn)
		}
	}
}
