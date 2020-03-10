package proliferate

import (
	"errors"
	"fmt"
	"time"

	"encoding/json"
)

// Block records are an interface so records can be abstract
type Block struct {
	ID           string                 `json:"_id"`
	Serial       int                    `json:"serial"`
	Timestamp    string                 `json:"timestamp"`
	Record       map[string]interface{} `json:"record"`
	Hash         string                 `json:"hash"`
	HashPrevious string                 `json:"hashPrevious"`
}

// Chain as a slice of blocks
type Chain []Block

//orderBlock returns record with new Block{} data from orderer
func (node *Node) orderBlock(record string, id string) Block {
	var raw map[string]interface{}

	if id == "" {
		id = NewID()
	}

	n := *node
	ts := time.Now()

	err := json.Unmarshal([]byte(record), &raw)
	if err != nil {
		n.Log(Message{
			Level: 1,
			Text:  err.Error(),
		})
	}

	lastBlock := n.Orderer.LastBlock(&n.Chain)

	serial := lastBlock.Serial + 1

	block := Block{
		ID:           id,
		Serial:       serial,
		Timestamp:    ts.String(),
		Record:       raw,
		HashPrevious: lastBlock.Hash,
	}

	block.Hash = Hash(block)

	return block
}

// LastBlock returns last block from storage or memory
func (node *Node) LastBlock() Block {
	n := *node

	// Just use last block in memory if storage disabled
	if n.Config.Couch.Enabled != true {
		return n.Chain[len(n.Chain)-1]
	}

	return n.LastBlockFromStorage()
}

//PushRecord calls PushBlock to push record as block to blockchain
func (node *Node) PushRecord(record string) {
	n := *node

	n.PushBlock(record, "")

	*node = n
}

//PushBlock pushes ordered block to the blockchain
func (node *Node) PushBlock(record string, id string) {
	n := *node

	if len(n.Chain) == 0 {
		n.Chain = append(n.Chain, n.Initialize())
	}

	block := n.orderBlock(record, id)

	n.Log(Message{
		Level: 5,
		Text:  "Pushing block: " + block.ID,
	})

	n.Chain = append(n.Chain, block)

	err := n.pushToStorage(block)
	if err != nil {
		n.Log(Message{
			Level: 1,
			Text:  err.Error(),
		})
	}

	n.EnforceMemoryLimit()

	*node = n

	// TODO EnforceMemoryLimit enable and test
}

// MarshalBlock returns block as string
func (node *Node) MarshalBlock(block Block) string {
	json, _ := json.Marshal(block)
	return string(json)
}

// Pushes block to chain on physical storage
func (node *Node) pushToStorage(block Block) error {
	n := *node
	c := n.Config.Couch

	if c.Enabled != true {
		return errors.New("Cannot push record: CouchDB is not enabled")
	}

	n.CouchPush(block)

	if n.DBExists() != true {
		return errors.New("CouchDB unavailable not running")
	}

	// TODO record interface logic should be considered here
	return nil
}

// BlockCount returns count of chains on block
func (node *Node) BlockCount() {
	//TODO count blocks on chain outside of blocks in memory
}

// EnforceMemoryLimit prevents blocks in memory to exceed limit
func (node *Node) EnforceMemoryLimit() {
	n := *node
	c := n.Config.Instance

	//if cap(n.Chain) > c.MemoryLimit && c.MemoryLimit != 0 {
	if len(n.Chain) > c.MemoryRecordLimit {
		n.Chain = append(n.Chain[:0], n.Chain[1:]...)
	}
	//fmt.Println("This is where it's breaking", c.MemoryRecordLimit)

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

	n.pushToStorage(block)

	return block
}

// ParseRecord returns Block.Record as string
func (node *Node) ParseRecord(block Block) string {
	// TODO properly parse this
	return fmt.Sprintf("%v", block.Record)
}
