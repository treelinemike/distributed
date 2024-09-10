package main

import (
	"log"
	"net/rpc"
)

func main() {
	log.Println("Starting client")

	// connect to host
	var err error
	remoteAccount, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Error connecting to banker")
	}

	// set name and balance
	var j int
	err = remoteAccount.Call("BankAccount.SetName", "Savings", &j)
	if err != nil {
		log.Fatal("Couldn't run set account name")
	}
	err = remoteAccount.Call("BankAccount.SetBalance", 100.24, &j)
	if err != nil {
		log.Fatal("Couldn't run set account balance")
	}

	// get account balance
	var bal float32
	err = remoteAccount.Call("BankAccount.GetBalance", j, &bal)
	if err != nil {
		log.Fatal("Couldn't run get account balance")
	}
	log.Printf("Balance is currently $%0.2f\n", bal)

	// print account info (note: this will print on )
	err = remoteAccount.Call("BankAccount.PrintAccount", j, &j)
	if err != nil {
		log.Fatal("Couldn't print account details")
	}

	// close connection to account
	remoteAccount.Close()

}
