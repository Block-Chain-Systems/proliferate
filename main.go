package main

import (
	"fmt"

	bc "./src/proliferate"
)

type Record struct {
	id    string
	value string
}

func main() {
	//chain := bc.Chain{}
	node := bc.Node{}

	node.PushBlock("{\"initial\":\"block\"}")
	node.PushBlock(Record{
		id:    bc.NewID(),
		value: "{\"test\":\"interface\"}",
	})

	fmt.Println(node.Chain[0].ID)
	fmt.Println(node.Chain[1].ID)
	fmt.Println(node.Chain[0].Hash)
	fmt.Println(node.Chain[1].Record)
	fmt.Println(node.Chain[1].HashPrevious)
	fmt.Println(node.Chain[1].Hash)

}
