package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

func main() {

	// create and register player API
	log.Println("Registering player API")
	thisapi := new(GFPlayerAPI)
	rpc.Register(thisapi)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("Error listening - is the server already running?")
		return
	}
	go http.Serve(l, nil)
	fmt.Println("Ready to play")

	// TODO: Find a better way to do this!
	gameover = false
	for !gameover {
		time.Sleep(100 * time.Microsecond) // TODO: pause needed to make sure we evaluate the gameover condition
	}
}
