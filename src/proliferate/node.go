package proliferate

type Node struct {
	Chain   Chain
	Config  Config
	Orderer Orderer
}

// Start starts the node
func (node *Node) Start() {
	n := *node

	n.Config = LoadConfig()

	*node = n
}
