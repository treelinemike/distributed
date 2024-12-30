package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

func main() {
	fid, err := os.Open("mazetest.json")
	if err != nil {
		log.Fatal("could not open json file")
	}
	jsondata, err := io.ReadAll(fid)
	if err != nil {
		log.Fatal("failed readall on USGS json")
	}
	fid.Close()

	var rstr Mazedata
	json.Unmarshal(jsondata, &rstr)

	// check that we have the correct number of elements in all data slices
	for _, e := range rstr.Elements {

		var expectedlen int64
		if e.Type < 100 {
			expectedlen = (2*rstr.M+1)*rstr.N + rstr.M
		} else {
			expectedlen = rstr.M * rstr.N
		}

		if len(e.Data) != int(expectedlen) {
			log.Fatalf("[type %d] incorrect data length: %d\n", e.Type, len(e.Data))
		}
	}

	log.Println(rstr.Title)
	log.Println(rstr.Author)
	log.Println(rstr.Description)
	log.Println(rstr.M)
	log.Println(rstr.N)
	log.Println(rstr.Elements)

}

// nested json structure for USGS river data
type Mazedata struct {
	Title       string        `json:"title"`
	Author      string        `json:"author"`
	Description string        `json:"description"`
	M           int64         `json:"m"`
	N           int64         `json:"n"`
	Elements    []Mazeelement `json:"elements"`
}

type Mazeelement struct {
	Type        int64     `json:"type"`
	Description string    `json:"description"`
	Data        []float64 `json:"data"`
}
