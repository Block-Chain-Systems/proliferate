package bc

import (
	"fmt"
	"time"

	"github.com/satori/go.uuid"
)

type Block struct {
	ID           string
	Index        int
	Timestamp    string
	Record       string
	Hash         string
	HashPrevious string
}

type Chain []Block

// prepareBlock evaluates chain and places record in new block
func (chain *Chain) prepareBlock(record string) Block {
	c := *chain
	ts := time.Now()

	prevBlock := c[len(c)-1]

	block := Block{
		ID:           NewID(),
		Index:        prevBlock.Index + 1,
		Timestamp:    ts.String(),
		Record:       record,
		HashPrevious: prevBlock.Hash,
	}

	block.Hash = Hash(block)

	return (block)
}

//pushBlock creates block and pushes record to chain
func (chain *Chain) PushBlock(record string) {
	c := *chain

	if len(c) == 0 {
		c = append(c, Initialize())
	}

	block := c.prepareBlock(record)
	c = append(c, block)
	*chain = c
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
