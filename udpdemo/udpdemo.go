// following: https://github.com/jeroendk/go-tcp-udp/blob/master/udpclient/main.go
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	log.Println("Starting program")
	serverAddr := "192.168.88.118:2390"
	/*udpAddr, err := net.ResolveUDPAddr("udp",serverAddr)
	if err != nil{
		log.Fatal("Failed to resolve UDP Address: ",err.Error())
	}*/

	udpnet, err := net.Dial("udp", serverAddr)
	if err != nil {
		log.Fatal("Failed to dial UDP Address: ", err.Error())
	}

	_, err = udpnet.Write([]byte("Hello, world!"))
	if err != nil {
		log.Fatal("Failed to write to UDP: ", err.Error())
	}

	// Read from the connection untill a new line is send
	data, err := bufio.NewReader(udpnet).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the data read from the connection to the terminal
	log.Print("Received: ", string(data))

}
