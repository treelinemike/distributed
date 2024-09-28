package main

import (
	"engg415/raft/common"
	"fmt"
	"log"
	"os"
)

func main() {

	// load cluster configuration
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run . configfile.yaml serverkey")
		return
	}
	filename := os.Args[1]
	selfkey := os.Args[2]
	servers := make(map[string]common.NetworkAddress)
	common.LoadRaftConfig(filename, servers)

	// make sure provided selfkey is in map from config file
	_, ok := servers[selfkey]
	if !ok {
		log.Fatal("Key ", selfkey, " is not in cluster config file")
	}
}
