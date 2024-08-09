package playingcards

import "strconv"

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


