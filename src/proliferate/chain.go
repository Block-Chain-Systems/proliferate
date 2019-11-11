package proliferate

import "time"

type Block struct {
	ID           string
	Serial       int
	Timestamp    string
	Record       interface{}
	Hash         string
	HashPrevious string
}

type Chain []Block

//prepareBlock returns record with new Block{} data from orderer
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

	return (block)
}

//pushBlock creates block and pushes record to chain
func (node *Node) PushBlock(record interface{}) {
	n := *node

	if len(n.Chain) == 0 {
		n.Chain = append(n.Chain, Initialize())
	}

	block := n.orderBlock(record)
	n.Chain = append(n.Chain, block)
	*node = n
}

//Initialize starts empty chain if no record exists
func Initialize() Block {
	ts := time.Now()

	block := Block{
		ID:        NewID(),
		Serial:    0,
		Timestamp: ts.String(),
		Record:    "{}",
	}

	block.Hash = Hash(block)

	return block
}
