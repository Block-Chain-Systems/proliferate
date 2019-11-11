package bc

import "fmt"

type Orderer struct{}

// IndexNext returns next available index int
func (orderer *Orderer) IndexNext(chain *Chain) int {

	// TODO Leader/Follower ordering

	c := *chain
	prevBlock := c[len(c)-1]
	return prevBlock.Index + 1
}
