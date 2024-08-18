package playingcards

import (
	"fmt"
	"math/rand"
)

type Deck struct {
	cards []Card
}

func (deck *Deck) AddCard(c Card) error {
	deck.cards = append(deck.cards, c)
	return nil
}

func (deck *Deck) Create() error {
	deck.cards = nil
    for i := 1; i <= 13; i++ {
		deck.AddCard(Card{i, Clubs})
		deck.AddCard(Card{i, Diamonds})
		deck.AddCard(Card{i, Hearts})
		deck.AddCard(Card{i, Spades})
	}
	return nil
}

// method to shuffle the deck
// shuffle (see: https://stackoverflow.com/questions/12264789/shuffle-array-in-go)
func (deck *Deck) Shuffle() error {
	for i := range deck.cards {
		j := rand.Intn(i + 1)
		deck.cards[i], deck.cards[j] = deck.cards[j], deck.cards[i]
	}
	return nil
}

func (deck Deck) Show() error {
	// let's check the shuffled deck
	space := ""
	for i := range deck.cards {
		fmt.Print(space, deck.cards[i].String())
		space = " "
	}
	fmt.Println()
	return nil
}

func (deck *Deck) TakeTopCard() Card {
	c := deck.cards[0]
	deck.cards = deck.cards[1:]
	return c
}

