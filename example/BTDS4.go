package main

import (
	"log"

	"github.com/darkautism/gds4"
)

func main() {
	log.Println("Connection to DS4")
	ds4, err := gds4.NewDS4("30:0E:D5:8E:7A:FC")
	if err != nil {
		if err.Error() == "host is down" {
			goto ReConn
		}
		log.Panic(err)
		return
	}
}
