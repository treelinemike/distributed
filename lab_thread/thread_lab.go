package main

import (
	"fmt"
	"sync"
	"time"
	//"net/http"
	//"io"
	//"bufio"
)

const num_threads = 8

var count int = 0
var countlock sync.Mutex
var wg sync.WaitGroup

func run_count() {
	defer wg.Done()
	for i := 0; i < 5000; i++ {
		//countlock.Lock()
		count += 1
        time.Sleep(0)
        count -= 1
		//countlock.Unlock()
	}
}

func main() {

	//resp,err := http.Get("https://www.gutenberg.org/cache/epub/100/pg100.tot")
	//defer resp.Body.Close()
	//fmt.Println("Attempted to open URL with error code = ",err)

	start_time := time.Now()
	wg.Add(num_threads)
	for i := 0; i < num_threads; i++ {
		go run_count()
	}

	wg.Wait()
	time_elapsed := time.Since(start_time).Seconds()
	fmt.Print("Count = ", count, " in ", time_elapsed, " sec\n")
}
