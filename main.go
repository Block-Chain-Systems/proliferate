package main

import (
	"fmt"

	bc "./src/proliferate"
)

func main() {
	//chain := bc.Chain{}
	node := bc.Node{}

	node.PushBlock("{\"Initial\":\"Block\"}")

	fmt.Println(node.Chain[0].ID)
	fmt.Println(node.Chain[1].ID)
	fmt.Println(node.Chain[0].Hash)
	fmt.Println(node.Chain[1].HashPrevious)
	fmt.Println(node.Chain[1].Hash)

}
