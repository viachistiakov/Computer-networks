package service

import (
	"fmt"
	"log"
	"sync"
)

type Node struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Segment struct {
	A Node `json:"a"`
	B Node `json:"b"`
}

type Message struct {
	SourceName string  `json:"sourceIP"`
	Segment    Segment `json:"segment"`
}

type Peer struct {
	Name     string
	IP       string
	Port     string
	NextIP   string
	NextPort string
	Logger   *log.Logger
	Stop     chan struct{}
	mu       sync.Mutex
	Segments []Segment
}

func (p *Peer) Address() string {
	return p.IP + ":" + p.Port
}

func (p *Peer) NextAddress() string {
	return p.NextIP + ":" + p.NextPort
}

func (p *Peer) StartPeer() {
	go p.listen()
	go p.say()
	<-p.Stop
	p.Logger.Printf("%s: peer shutdown\n", p.Address())
	fmt.Println("Peer shutting down...")
}
