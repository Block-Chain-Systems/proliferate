package proliferate

import (
	"fmt"
	"time"
)

// Block records are an interface so records can be abstract
type Block struct {
	ID           string      `json:"id"`
	Serial       int         `json:"serial"`
	Timestamp    string      `json:"timestamp"`
	Record       interface{} `json:"record"`
	Hash         string      `json:"hash"`
	HashPrevious string      `json:"hashPrevious"`
}

// Chain as a slice of blocks
type Chain []Block

//orderBlock returns record with new Block{} data from orderer
func (node *Node) orderBlock(record interface{}) Block {
	n := *node
	ts := time.Now()

	lastBlock := n.Orderer.LastBlock(&n.Chain)

	block := Block{
		ID:           NewID(),
		Serial:       node.Orderer.SerialNext(&n.Chain),
		Timestamp:    ts.String(),
		Record:       record,
		HashPrevious: lastBlock.Hash,
	}

	block.Hash = Hash(block)

	return block
}

//PushBlock pushes ordered block to the blockchain
func (node *Node) PushBlock(record interface{}) {
	n := *node

	if len(n.Chain) == 0 {
		n.Chain = append(n.Chain, n.Initialize())
	}

	block := n.orderBlock(record)

	n.Log(Message{
		Level: 5,
		Text:  "Pushing block: " + block.ID,
	})

	n.Chain = append(n.Chain, block)

	n.pushToStorage(block)

	*node = n

	// TODO EnforceMemoryLimit enable and test
}

// Pushes block to chain on physical storage
func (node *Node) pushToStorage(block Block) {
	n := *node
	//c := n.Config

	if n.DBExists() != true {
		n.Log(Message{
			Level: 2,
			Text:  "Cannot find couchDB database or couchDB not running",
		})
	}

	// TODO record interface logic should be called here

	//if c.Couch.Enabled == true {
	//	n.CouchPut(fmt.Sprintf("%v", block.Record))
	//}
}

// BlockCount returns count of chains on block
func (node *Node) BlockCount() {
	//TODO count blocks on chain outside of blocks in memory
}

// EnforceMemoryLimit prevents blocks in memory to exceed limit
func (node *Node) EnforceMemoryLimit() {
	n := *node
	c := n.Config.Instance

	if cap(n.Chain) > c.MemoryLimit && c.MemoryLimit != 0 {
		n.Chain = append(n.Chain[:0], n.Chain[1:]...)
	}

	*node = n
}

//Initialize starts empty chain if no blocks exist
func (node *Node) Initialize() Block {
	n := *node
	ts := time.Now()

	n.Log(Message{
		Level: 3,
		Text:  "No blocks found. Creating genesis block",
	})

	block := Block{
		ID:        NewID(),
		Serial:    0,
		Timestamp: ts.String(),
	}

	block.Hash = Hash(block)

	return block
}

// TODO properly parse Block.Record
func (node *Node) ParseRecord(block Block) string {
	return fmt.Sprintf("%v", block.Record)
}
