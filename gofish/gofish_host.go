package main

import (
	"engg415/playingcards"
	"fmt"
)

func main() {

	// create and shuffle a standard deck
	deck := new(playingcards.Deck)
	deck.Create()
	deck.Create()
    deck.Shuffle()
	fmt.Println("Shuffled deck:")
	deck.Show()
	fmt.Println("Top card is: ", deck.TakeTopCard().String())
	fmt.Println("Deck is now:")
	deck.Show()
	deck.Shuffle()
	fmt.Println("After shuffling again:")
	deck.Show()
}
