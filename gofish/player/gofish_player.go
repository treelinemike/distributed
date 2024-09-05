package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

func main() {

	// format log
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

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
	log.Println("Ready to play")

	// TODO: Find a better way to do this!
	gameover = false
	for !gameover {
		time.Sleep(10 * time.Millisecond) // TODO: pause needed to make sure we evaluate the gameover condition
	}
}
