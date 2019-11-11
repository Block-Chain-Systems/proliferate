package proliferate

//import "fmt"

type Orderer struct{}

// TODO Leader/Follower peer ordering

// IndexNext returns next available index int
func (orderer *Orderer) IndexNext(chain *Chain) int {
	c := *chain
	//prevBlock := c[len(c)-1]
	prevBlock := orderer.LastBlock(&c)
	return prevBlock.Index + 1
}

func (orderer *Orderer) LastBlock(chain *Chain) Block {
	c := *chain
	return c[len(c)-1]
}
