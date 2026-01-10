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

type WrappedResponse struct {
	r   common.RespToClient
	err error
}

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
	allSvrKeys := make([]string, 0, len(servers))
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

		// also add this to a slice of server keys
		allSvrKeys = append(allSvrKeys, svr_key)
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

	// open the input text file
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

	// open the ground truth output text file
	fileout, err := os.OpenFile("wordlist_sent.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("Error opening output file")
	}
	defer fileout.Close()
	writer := bufio.NewWriter(fileout)
	defer writer.Flush()

	// select random strings from the slice to store in the raft cluster
	for range 100 {

		// randomly select our next word and remove it from the list
		idx := rng.IntN(len(strings))
		thisWord := strings[idx]
		log.Printf("Choosing word %d of %d: %s\n", idx, len(strings), thisWord)
		strings = slices.Delete(strings, idx, idx+1)

		svrChoice := ""
		committed := false
		for !committed {

			// randomly select a server to receive this word
			// making sure that we *think* the server is currently active

			for svrChoice == "" {
				svrChoiceTemp := allSvrKeys[rand.IntN(len(allSvrKeys))]
				if isactive[svrChoiceTemp] {
					svrChoice = svrChoiceTemp
				}
			}
			log.Printf("Attempting to send to server %v\n", svrChoice)

			// send this word to server
			serverResp := make(chan WrappedResponse)
			go func(c chan WrappedResponse, word string) {
				// TODO: ACTUALLY SUBMIT HERE VIA RPC
				var wr WrappedResponse
				var r common.RespToClient
				err := servers[svrChoice].Handle.Call("RaftAPI.ProcessClientRequest", []string{thisWord}, &r)
				wr.err = err
				wr.r = r
				c <- wr
			}(serverResp, thisWord)
			var resp WrappedResponse
			select {
			case resp = <-serverResp:
			case <-time.After(200 * time.Millisecond):
				// remove this server from active list and try to reconnect
				log.Printf("Timeout submitting word to server %v\n", svrChoice)
				isactive[svrChoice] = false
				go reconnectToHost(svrChoice)
				svrChoice = ""
				continue // re-attempt!
			}

			if resp.err != nil {
				// remove this server from active list and try to reconnect
				log.Printf("Error submitting word to server: %v\n", resp.err)
				isactive[svrChoice] = false
				go reconnectToHost(svrChoice)
				svrChoice = ""
				continue // re-attempt!
			}

			if !resp.r.Committed {
				log.Printf("Server %v is not leader, redirecting to server %v\n", svrChoice, resp.r.LeaderID)
				svrChoice = resp.r.LeaderID
				continue // re-attempt!
			}

			// word has been committed
			committed = true

			// write this word to the ground truth output file
			_, err := writer.WriteString(thisWord + "\n")
			if err != nil {
				log.Fatalf("Error writing '%v' to output file\n", thisWord)
			}

			// wait for a second
			time.Sleep(1 * time.Second)
		}
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
