package main

//import "errors"
import "fmt"
import "strconv"
import "math/rand"

type Suit int
const (
   Clubs Suit = iota
   Diamonds
   Hearts
   Spades
   )

type Card struct {
    val int
    suit Suit
}

func (c Card) String() string{
    var str string
    
    if(c.val > 1 && c.val < 11){
        str += strconv.Itoa(c.val)
    } else {
        switch c.val{
        case 1:
            str += "A"
        case 11:
            str += "J"
        case 12:
            str += "Q"
        case 13:
            str += "K"
        }
    }
    switch c.suit{
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


// DECK MANAGEMENT

type Deck struct{
    cards []Card
}

func (deck *Deck) AddCard(c Card) error{
    deck.cards = append(deck.cards,c)
    return nil
}

func (deck *Deck) Create() error{
    for i:=1; i<=13; i++ {
        deck.AddCard(Card{i,Clubs})
        deck.AddCard(Card{i,Diamonds})
        deck.AddCard(Card{i,Hearts})
        deck.AddCard(Card{i,Spades})
    }
    return nil
}
// function to shuffle the full deck
// shuffle (see: https://stackoverflow.com/questions/12264789/shuffle-array-in-go)
func (deck *Deck) Shuffle() error{
    for i:= range deck.cards {
        j := rand.Intn(i+1)
       deck.cards[i], deck.cards[j] = deck.cards[j], deck.cards[i]

    }
    return nil
}

func (deck *Deck) TakeTopCard() Card{

    c := deck.cards[0]

    // TODO: pop off top card here

    return c
}


// the "Hand" structre holding info about the current hand
// must have a capital "H" indicating that it is exported
type Hand struct{
    numCards uint // initialized to zero!
}

// method to update number of cards in hand
func (h *Hand) AddCards( numCardsToAdd uint, reply *int) error {
    h.numCards += numCardsToAdd
    return nil
}

func main() {
    fmt.Print("Hello! \u2663\u2666\u2665\u2660 \n")

    //hand := new(Hand)

    // create and shuffle a standard deck
    deck := new(Deck)
    deck.Create()
    deck.Shuffle();

    // let's check the shuffled deck
    fmt.Println("Shuffled deck:");
    space := ""
    for i:= range deck.cards{
        fmt.Print(space,deck.cards[i].String())
        space = " "
    }
    fmt.Println()

    
    fmt.Println("Top card is: ",deck.TakeTopCard().String())
}
