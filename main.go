package main

import (
	//"fmt"

	bc "./src/proliferate"
)

type record struct {
	id    string
	value string
	test  []string
}

func main() {
	//chain := bc.Chain{}
	node := bc.Node{}

	node.Start()
	node.PushBlock(`{"initial":"block"}`)

	testBlock := `{"Your":"Mom"}`
	node.PushBlock(testBlock)

	bc.DumpChain(node.Chain)
	//node.VerifyIdentity()
}
