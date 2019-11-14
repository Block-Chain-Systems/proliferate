package proliferate

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

// Orderer struct
type Orderer struct{}

// TODO Leader/Follower peer ordering

// SerialNext returns next available index int
func (orderer *Orderer) SerialNext(chain *Chain) int {
	c := *chain
	prevBlock := orderer.LastBlock(&c)
	return prevBlock.Serial + 1
}

// LastBlock returns last block on chain
func (orderer *Orderer) LastBlock(chain *Chain) Block {
	// TODO last block should check consensus

	c := *chain
	return c[len(c)-1]
}

// ReplicateBlock pushes block to peers and awaits consensus
func (orderer *Orderer) ReplicateBlock() {

}

//NewID generates UUID V4 ID
func NewID() string {
	id := uuid.Must(uuid.NewV4())
	return fmt.Sprintf("%s", id)
}
