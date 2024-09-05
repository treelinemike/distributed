package playingcards

import (
	"log"
	"strconv"
)

type Suit int

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
)

// values in Card need to be exported to use as argument in RPC
type Card struct {
	Val      int
	CardSuit Suit
}

type TestCard struct {
	Val, Suit int
}

type MyArgType struct {
	X, Y     int
	ThisName string
	ThisSuit Suit
}

func NumToCardChar(val int) string {
	if val > 1 && val < 11 {
		return strconv.Itoa(val)
	} else {
		switch val {
		case 1:
			return "A"
		case 11:
			return "J"
		case 12:
			return "Q"
		case 13:
			return "K"
		default:
			log.Panic("Invalid numerical card rank!")
			return ""
		}
	}
}

func (c Card) String() string {
	var str string
	str += NumToCardChar(c.Val)
	switch c.CardSuit {
	case Clubs:
		str += "\u2663"
	case Diamonds:
		str += "\u2666"
	case Hearts:
		str += "\u2665"
	case Spades:
		str += "\u2660"
	}
	return str
}
