package playingcards

import (
	"fmt"
	"math/rand"
    "errors"
)

type Deck struct {
    cards []Card
}

// wrapper for AddCard() that statisfies the RPC interface
func (deck *Deck) AddCardRPC(card Card,resp *int) error{
    err := deck.AddCard(card)
    if( err != nil ){
        return errors.New("Could not add card")
    }
    fmt.Println("We got a card! Current hand:")
    deck.Show()
    return nil
}

func (deck *Deck) ResetDeckRPC(args int, resp *int) error{
    fmt.Println("Resetting deck")
    deck.cards = nil
    return nil
}

func (deck *Deck) AddCard(c Card) error {
	if( c.Val < 1 || c.Val > 13){
        return errors.New("Invalid card")
    }
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

func (deck Deck) NumCards() int{
    return len(deck.cards)
}

func (deck Deck) Show() error {
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

