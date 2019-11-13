package proliferate

import (
	//"fmt"

	"encoding/json"
	"os"
)

type Node struct {
	ID      string
	Detail  PeerDetail
	Chain   Chain
	Config  Config
	Orderer Orderer
}

type PeerDetail struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Organization string `json:"organization"`
}

// Start starts the node
func (node *Node) Start() {
	n := *node

	n.ParseIdentity()
	n.Config = LoadConfig()

	n.Log(Message{
		Level: 4,
		Text:  "Network started",
	})

	*node = n
}

// LoadJSON returns json as struct (TODO!)
func (node *Node) ParseIdentity() {
	n := *node
	file := "network/node/id.json"
	var detail PeerDetail

	n.Log(Message{
		Level: 4,
		Text:  "Parsing node identity in network/id",
	})

	identityFile, err := os.Open(file)
	defer identityFile.Close()

	if err != nil {
		LogRaw(Message{
			Level: 1,
			Text:  err.Error(),
		})
	}
	jsonParser := json.NewDecoder(identityFile)
	jsonParser.Decode(&detail)

	n.Detail = detail
	n.ID = detail.ID

	*node = n
}
