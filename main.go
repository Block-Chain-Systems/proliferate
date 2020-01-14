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
	node.PushRecord(`{"initial":"block"}`)

	testBlock := `{"Your":"Mom"}`
	node.PushRecord(testBlock)

	bc.DumpChain(node.Chain)

	//fmt.Println(node.VerifyLastBlock())
	//fmt.Println(node.LoadDocumentsFromStorage())

	//fmt.Println(node.LoadIDsFromStorage)
	node.LoadChainFromStorage()

	node.LastBlockFromStorage()
	//node.VerifyIdentity()
	//fmt.Println(node.LoadIDsFromStorage())

	/*
		testRecord := make(map[string]interface{})
		testRecord["Your"] = "Mom"

		verifyBlock := bc.Block{
			ID:           "24646c5a-823b-47d8-9a30-9f51e02973cd",
			Serial:       1,
			Timestamp:    "2019-12-12 16:15:22.018413053 -0700 MST m=+0.010098811",
			Record:       testRecord,
			Hash:         "45c7f240532cbe6326af1406a59ef125b47f266b11466ca6ca45d75ca312d120",
			HashPrevious: "1c7f21e72e8a783bb8805c013ab304ae90061bcff447b77dbae5cd0550fde671",
		}
		fmt.Println(bc.VerifyHash(verifyBlock, verifyBlock.Hash))
	*/
}
