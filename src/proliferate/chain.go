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
func (node *Node) orderBlock(record string) Block {
	var raw map[string]interface{}

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

	block := Block{
		ID:           NewID(),
		Serial:       node.Orderer.SerialNext(&n.Chain),
		Timestamp:    ts.String(),
		Record:       raw,
		HashPrevious: lastBlock.Hash,
	}

	block.Hash = Hash(block)

	return block
}

//PushBlock pushes ordered block to the blockchain
func (node *Node) PushBlock(record string) {
	n := *node

	if len(n.Chain) == 0 {
		n.Chain = append(n.Chain, n.Initialize())
	}

	block := n.orderBlock(record)

	//fmt.Println("\n\n\n----PrepareBlock----")
	//fmt.Println(n.PrepareBlock(block))
	//fmt.Println("----/PrepareBlock----\n\n\n")

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

	*node = n

	// TODO EnforceMemoryLimit enable and test
}

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
	} else {
		n.CouchPush(block)
	}

	if n.DBExists() != true {
		return errors.New("CouchDB unavailable not running")
	}

	// TODO record interface logic should be considered here
	return nil
}

//func (node *Node) PrepareBlock(block Block) string {
//
//}

// BlockCount returns count of chains on block
func (node *Node) BlockCount() {
	//TODO count blocks on chain outside of blocks in memory
}

// EnforceMemoryLimit prevents blocks in memory to exceed limit
func (node *Node) EnforceMemoryLimit() {
	n := *node
	c := n.Config.Instance

	if cap(n.Chain) > c.MemoryLimit && c.MemoryLimit != 0 {
		n.Chain = append(n.Chain[:0], n.Chain[1:]...)
	}

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

// TODO properly parse Block.Record
func (node *Node) ParseRecord(block Block) string {
	return fmt.Sprintf("%v", block.Record)
}
