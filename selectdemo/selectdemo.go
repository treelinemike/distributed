package main

import (
	"log"
	"time"
)

// main function

func main() {

	// make a channel
	votechan := make(chan bool)

	ballotsReceived := 1 // vote for self!
	votecount := 1       // vote for self!

	for i := 0; i < 3; i++ {
		go func(votechannel chan bool) {
			// simulate voting process

			time.Sleep(100 * time.Microsecond)

			// send vote result to channel
			votechannel <- true
		}(votechan)
	}

	mytimer := time.NewTimer(2 * time.Second)
	mytimer.Stop()
	mytimer.Reset(2 * time.Second)

	for done := false; !done; {
		select {
		case vote := <-votechan:

			// log that we got a ballot
			ballotsReceived += 1

			// determne whether we got the vote
			if vote {
				votecount += 1
				println("Received a vote")
			} else {
				println("Did not receive a vote")
			}

			// exit if we have a majority
			if votecount >= 2 {
				println("Received majority votes, winning election")
				// change to leader status here
				done = true
			}

			// exit if we have received all ballots
			if ballotsReceived == 3 {
				println("Received all ballots, election over")
				done = true
			}

		// timeout case
		case <-mytimer.C:
			log.Println("Timeout waiting for vote")
			done = true
		}
	}

	log.Printf("Total votes received at exiting time: %d", votecount)

}
