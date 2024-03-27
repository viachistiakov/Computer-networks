package main

import (
	"fmt"
	"log"
	"os"
	"peer/src/service"
)

func main() {
	var name, peerIP, port, nextPeerIP, nextPort string

	fmt.Print("Enter peer name: ")
	fmt.Scan(&name)
	fmt.Print("Enter peer IP: ")
	fmt.Scan(&peerIP)
	fmt.Print("Enter peer port: ")
	fmt.Scan(&port)
	fmt.Print("Enter next peer IP: ")
	fmt.Scan(&nextPeerIP)
	fmt.Print("Enter next peer port: ")
	fmt.Scan(&nextPort)

	logger := log.Default()
	f, err := os.OpenFile("logFile.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	logger.SetOutput(f)

	p := &service.Peer{
		Name:     name,
		IP:       peerIP,
		Port:     port,
		NextIP:   nextPeerIP,
		NextPort: nextPort,
		Logger:   logger,
		Stop:     make(chan struct{}),
		Segments: make([]service.Segment, 0),
	}

	p.StartPeer()
}
