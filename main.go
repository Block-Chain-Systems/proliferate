package main

import (
	"fmt"

	"./src/bc"
)

func main() {
	chain := bc.Chain{}

	chain.PushBlock("{\"Initial\":\"Block\"}")

	fmt.Println(chain[0].ID)
	fmt.Println(chain[1].ID)
	fmt.Println(chain[1].HashPrevious)

	//fmt.Println(chain)
}
