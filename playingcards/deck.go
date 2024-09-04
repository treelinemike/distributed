package playingcards

import (
	"errors"
	"math/rand"
	"time"
)

type Deck struct {
	cards []Card
}

func (deck *Deck) AddCard(c Card) error {
	if c.Val < 1 || c.Val > 13 {
		return errors.New("invalid card")
	}
	deck.cards = append([]Card{c}, deck.cards...)
	return nil
}

func (deck *Deck) Reset() error {
	deck.cards = nil
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
// https://stackoverflow.com/questions/12321133/how-to-properly-seed-random-number-generator
func (deck *Deck) Shuffle() error {

	// TODO: check whether we can get consistent performance with a specified seed
	// if so could use this for grading purposes
	rand.Seed(time.Now().UnixNano())

	// shuffle cards
	for i := range deck.cards {
		j := rand.Intn(i + 1)
		deck.cards[i], deck.cards[j] = deck.cards[j], deck.cards[i]
	}
	return nil
}

func (deck Deck) NumCards() int {
	return len(deck.cards)
}

func (deck Deck) String() (s string) {
	space := ""
	for i := range deck.cards {
		s = s + space + deck.cards[i].String()
		space = " "
	}
	return
}

func (deck *Deck) TakeTopCard() Card {
	c := deck.cards[0]
	deck.cards = deck.cards[1:]
	return c
}
