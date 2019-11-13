package proliferate

import (
	"fmt"
	"time"
)

type Block struct {
	ID           string
	Serial       int
	Timestamp    string
	Record       interface{}
	Hash         string
	HashPrevious string
}

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

//pushBlock creates block and pushes record to chain
func (node *Node) PushBlock(record interface{}) {
	n := *node

	if len(n.Chain) == 0 {
		n.Chain = append(n.Chain, n.Initialize())
	}

	block := n.orderBlock(record)

	n.Log(Message{
		Level: 5,
		Text:  fmt.Sprintf("Pushing block: %v", block.ID),
	})

	n.Chain = append(n.Chain, block)
	*node = n
}

//Initialize starts empty chain if no record exists
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
