package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"slices"
)

func main() {

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
	source := rand.New(rand.NewPCG(102, 3145))
	rng := rand.New(source)

	// select random strings from the slice to store in the raft cluster
	for range 100 {
		idx := rng.IntN(len(strings))
		fmt.Printf("%d (of %d): %s\n", idx, len(strings), strings[idx])
		strings = slices.Delete(strings, idx, idx+1)
	}
}
