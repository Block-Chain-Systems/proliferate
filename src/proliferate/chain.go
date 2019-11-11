package proliferate

import (
	"fmt"
	"time"

	"github.com/satori/go.uuid"
)

var order = Orderer{}

func (node *Node) prepareBlock(record string) Block {
	n := *node
	ts := time.Now()

	lastBlock := order.LastBlock(&n.Chain)

	block := Block{
		ID:           NewID(),
		Index:        node.Orderer.IndexNext(&n.Chain),
		Timestamp:    ts.String(),
		Record:       record,
		HashPrevious: lastBlock.Hash,
	}

	block.Hash = Hash(block)

	return (block)
}

//pushBlock creates block and pushes record to chain
func (node *Node) PushBlock(record string) {
	n := *node

	if len(n.Chain) == 0 {
		n.Chain = append(n.Chain, Initialize())
	}

	block := n.prepareBlock(record)
	n.Chain = append(n.Chain, block)
	*node = n
}

//Initialize initialized chain if no record exists
func Initialize() Block {
	ts := time.Now()

	block := Block{
		ID:        NewID(),
		Index:     0,
		Timestamp: ts.String(),
		Record:    "{}",
	}

	block.Hash = Hash(block)

	return block
}

//NewID generates UUID V4 ID
func NewID() string {
	id := uuid.Must(uuid.NewV4())
	return fmt.Sprintf("%s", id)
}
