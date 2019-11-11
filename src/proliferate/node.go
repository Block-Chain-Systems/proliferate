package proliferate

type Node struct {
	Chain   Chain
	Orderer Orderer
}

type Block struct {
	ID           string
	Serial       int
	Index        int
	Timestamp    string
	Record       string
	Hash         string
	HashPrevious string
}

type Chain []Block
