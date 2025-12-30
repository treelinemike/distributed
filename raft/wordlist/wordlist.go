package main

import (
	"bufio"
	"log"
	"os"
	"slices"
)

func main() {

	// open the text file
	file, err := os.OpenFile("google-10000-english-usa-no-swears.txt", os.O_RDONLY, os.ModePerm)
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

	// open the text file
	file2, err := os.OpenFile("remove.txt", os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	defer file2.Close()

	// load each line of file into an element of a slice
	var badstrings []string
	scanner2 := bufio.NewScanner(file2)
	for scanner2.Scan() {
		//log.Printf("Read string: %s\n", scanner2.Text())
		badstrings = append(badstrings, scanner2.Text())
	}

	// remove bad words
	for _, badword := range badstrings {
		log.Printf("Removing occurrences of bad word: %s\n", badword)
		for idx, thisword := range strings {
			if thisword == badword {
				log.Printf("Removing bad word: %s\n", thisword)
				strings = slices.Delete(strings, idx, idx+1)
			}
		}
	}

	// sort the strings slice alphabetically
	slices.Sort(strings)

	// save the cleaned list
	outputfile, err := os.OpenFile("wordlist.txt", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer outputfile.Close()

	writer := bufio.NewWriter(outputfile)
	for _, thisword := range strings {
		_, err := writer.WriteString(thisword + "\n")
		if err != nil {
			log.Fatalf("Error writing to output file: %v", err)
		}
	}
	writer.Flush()

	/*
		// seed the RNG so we can reproduce results
		source := rand.New(rand.NewPCG(102, 9945))
		rng := rand.New(source)

		// select random strings from the slice to store in the raft cluster
		for range 100 {
			idx := rng.IntN(len(strings))
			fmt.Printf("%d (of %d): %s\n", idx, len(strings), strings[idx])
			strings = slices.Delete(strings, idx, idx+1)
		}
	*/
}
