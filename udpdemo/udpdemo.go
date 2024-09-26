// following: https://github.com/jeroendk/go-tcp-udp/blob/master/udpclient/main.go
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
)

func main() {
	log.Println("Starting program")

	// see if we got an IP as a command line argument
	serverAddr := "192.168.88.201:2390"
	if len(os.Args) == 2 {
		serverAddr = os.Args[1] + ":2390"
	}

	/*udpAddr, err := net.ResolveUDPAddr("udp",serverAddr)
	if err != nil{
		log.Fatal("Failed to resolve UDP Address: ",err.Error())
	}*/

	udpnet, err := net.Dial("udp", serverAddr)
	if err != nil {
		log.Fatal("Failed to dial UDP Address: ", err.Error())
	}

	lightstatus := false
	for {
		if !lightstatus {
			_, err = udpnet.Write([]byte("lighton"))
			lightstatus = true

			_, err := exec.Command("blink1-on.sh").Output()
			if err != nil {
				log.Println("Could not turn on blink(1) indicator")
			}

		} else {
			_, err = udpnet.Write([]byte("lightoff"))
			lightstatus = false

			_, err := exec.Command("blink1-off.sh").Output()
			if err != nil {
				log.Println("Could not turn on blink(1) indicator")
			}

		}
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
		time.Sleep(time.Second * 1)
	}

}
