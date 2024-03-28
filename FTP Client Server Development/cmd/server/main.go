package main

import (
	"lab6_1/server"
	"log"
)

func main() {
	conf := &server.Config{
		Root: "root",
		Host: "0.0.0.0",
		Port: 55321,
		User: "user",
		Pass: "1234test",
	}
	s := server.NewServer(conf)
	err := s.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
