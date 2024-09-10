package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type BankAccount struct {
	name    string
	balance float32
}

func (ba *BankAccount) SetName(newName string, resp *int) error {
	ba.name = newName
	return nil
}

func (ba *BankAccount) SetBalance(newBalance float32, resp *int) error {
	ba.balance = newBalance
	return nil
}

func (ba *BankAccount) GetBalance(args int, resp *float32) error {
	*resp = ba.balance
	return nil
}

func (ba *BankAccount) PrintAccount(args int, resp *int) error {
	log.Printf("Account: %s; Balance $%0.2f\n", ba.name, ba.balance)
	return nil
}

func main() {

	// show that we've started the server
	log.Println("Starting banker")

	// create a new bank account
	yourAccount := new(BankAccount)

	// allow remote access to this new bank account
	rpc.Register(yourAccount)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("Error listening - is the server already running?")
		return
	}
	go http.Serve(l, nil)
	log.Println("Serving bank account")

	// spin while allowing connections
	for {

	}

}
