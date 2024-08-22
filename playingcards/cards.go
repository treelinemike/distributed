package playingcards

import "strconv"
import "fmt"

type Suit int
const (
   Clubs Suit = iota
   Diamonds
   Hearts
   Spades
   )

// values in Card need to be exported to use as argument in RPC
type Card struct {
    Val int
    CardSuit Suit
}

type TestCard struct{
    Val,Suit int
}

type MyArgType struct{
    X,Y int
    ThisName string
    ThisSuit Suit
}

func (tc *TestCard) TestRPC(a MyArgType, b *int) error{
    fmt.Println("Executing RPC. X = ",a.X,a.ThisSuit)
    return nil
}

func (c Card) String() string{
    var str string
    
    if(c.Val > 1 && c.Val < 11){
        str += strconv.Itoa(c.Val)
    } else {
        switch c.Val{
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
    switch c.CardSuit{
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


