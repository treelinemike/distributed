package main

import (
	"bufio"
	"engg415/raft/common"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/rpc"
	"os"
	"slices"
	"sync"
	"time"
)

var servers map[string]common.NetworkAddress
var isactive map[string]bool

// repeatedly attempt to reconnect to host
func reconnectToHost(svr string) {
	log.Printf("Attempting to reconnect to server %v (%s:%s)\n", svr, servers[svr].Address, servers[svr].Port)
	for {
		host, err := rpc.DialHTTP("tcp", servers[svr].Address+":"+servers[svr].Port)
		if err == nil {
			log.Printf("Reconnected to server %v (%s:%s)\n", svr, servers[svr].Address, servers[svr].Port)

			// store handle to server in struct in map
			// Go doesn't make this easy
			tempsvrdata := servers[svr]
			tempsvrdata.Handle = host
			servers[svr] = tempsvrdata

			// update to indicate that server is back online
			isactive[svr] = true

			// break out of loop
			break
		}
	}
}

func main() {

	var err error

	// format log time in microseconds
	// will start logging to file as soon as we know which host key we have been assigned
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	// load cluster configuration
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run . configfile.yaml")
		return
	}
	filename := os.Args[1]
	servers = make(map[string]common.NetworkAddress)
	isactive = make(map[string]bool)

	// load config
	log.Println("Loading cluster configuration file")
	var jsonfilebase string
	var t common.Timeout
	common.LoadRaftConfig(filename, servers, &t, &jsonfilebase)

	// configure logging
	logfilename := "log_client.txt"
	log.Printf("Creating log file: %s", logfilename)
	logfid, err := os.OpenFile(logfilename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Panicf("Could not create log file: %v\n", err)
	}
	defer logfid.Close()
	mux := io.MultiWriter(os.Stdout, logfid)
	log.SetOutput(mux)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds) // for good measure?

	// connect to each server in the cluster
	// use a bunch of goroutines with timeouts
	// ok to block with a wait group here beause we do actually need to connect to all of our servers
	// TODO: better handle the case when a server is initially crashed?
	log.Println("Trying to connect to all servers now...")
	var wg sync.WaitGroup
	for svr_key, svr_data := range servers { // ranging over a map returns (key, value) pairs

		wg.Add(1)
		go func(svr_key string) { // note: we can still access servers due to scoping of goroutine
			defer wg.Done()
			// try connect to host repeatedly (even if OS times out the attempt very quickly)
			// we will use our own timeout on the connection attempt
			for {
				//log.Printf("Attempting connection to host %s:%s\n", servers[svr].Address, servers[svr].Port)
				host, err := rpc.DialHTTP("tcp", servers[svr_key].Address+":"+servers[svr_key].Port)
				if err == nil {
					log.Printf("Connected to host %s:%s\n", servers[svr_key].Address, servers[svr_key].Port)

					// store handle to server in struct in map
					// Go doesn't make this easy
					tempsvrdata := svr_data
					tempsvrdata.Handle = host
					servers[svr_key] = tempsvrdata

					// set "up" status
					isactive[svr_key] = true

					// break out of loop
					break
				}
			}
		}(svr_key)

	}
	c := make(chan struct{})
	go func() {
		wg.Wait()
		defer close(c)
	}()
	select {
	case <-c:
		log.Println("Connected to all servers")
	case <-time.After(5 * time.Minute):
		log.Fatal("Timeout connecting to servers")
	}

	// open the text file
	file, err := os.OpenFile("wordlist.txt", os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	defer file.Close()

	// load each line of file into an element of a slice
	var strings []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		strings = append(strings, scanner.Text())
	}

	// seed the RNG so we can reproduce results
	source := rand.New(rand.NewPCG(199, 3135))
	rng := rand.New(source)

	// select random strings from the slice to store in the raft cluster
	for range 100 {
		idx := rng.IntN(len(strings))
		fmt.Printf("%d (of %d): %s\n", idx, len(strings), strings[idx])
		strings = slices.Delete(strings, idx, idx+1)
	}

	// stop all servers
	for svr := range servers {
		log.Printf("Attempting to stop server %v\n", svr)
		var temp int
		err := servers[svr].Handle.Call("RaftAPI.StopServer", temp, &temp)

		if err != nil {
			log.Printf("Could not stop server %s: %v\n", svr, err)
		}
	}
}
