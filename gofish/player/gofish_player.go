package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"time"
)

func main() {

	// format log
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// if there is one command line argument, assume it is a port
	// and try to connect on it
	var port string
	switch len(os.Args) {
	case 1:
		port = ":1234"
	case 2:
		_, err := strconv.ParseInt(os.Args[1], 10, 16)
		if err != nil {
			log.Fatal("Invalid port")
		} else {
			port = ":" + os.Args[1]
		}
	default:
		log.Fatal("Too many arguments provided")
	}

	// create and register player API
	log.Println("Registering player API for access on port", port[1:])
	thisapi := new(GFPlayerAPI)
	rpc.Register(thisapi)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", port)
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
