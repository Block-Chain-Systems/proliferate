package main

import (
	"fmt"

	bc "./src/proliferate"
)

type Record struct {
	id    string
	value string
	test  []string
}

func main() {
	//chain := bc.Chain{}
	node := bc.Node{}

	node.Start()
	node.PushBlock("{\"initial\":\"block\"}")

	node.PushBlock(Record{
		id:    bc.NewID(),
		value: "{\"test\":\"interface\"}",
		test:  []string{0: "one", 1: "two", 2: "three"},
	})

	//bc.DumpChain(node.Chain)

	//node.CouchTest()

	fmt.Println(bc.NewID())
}
