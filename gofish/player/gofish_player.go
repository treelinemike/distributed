package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
)

func main() {

	fmt.Println("Starting player")

	// create and register player API
	thisapi := new(GFPlayerAPI)
	rpc.Register(thisapi)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		fmt.Println("Error listening - is the server already running?")
		return
	}
	go http.Serve(l, nil)
	fmt.Println("Ready to play!")

	// TODO: Find a better way to do this!
	for {
	}
}
