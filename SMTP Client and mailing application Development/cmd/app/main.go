package main

import (
	"lab6_2/internal/client"
	"lab6_2/internal/repository"
	"log"
)

func main() {
	db, err := repository.NewRepository("iu9networkslabs:Je2dTYr6@tcp(students.yss.su)/iu9networkslabs")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	c := client.NewClient(db, "dts21@dactyl.su", "12345678990DactylSUDTS",
		"mail.nic.ru", 465)
	c.SendMailList()
}
