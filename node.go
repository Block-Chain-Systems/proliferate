package proliferate

import (
	"encoding/json"
	//"fmt"
	"os"
	"path"
)

// Node struct
type Node struct {
	ID      string
	Detail  PeerDetail
	Chain   Chain `json:"docs"`
	Config  Config
	Orderer Orderer
	Peers   []PeerDetail
	member  member
}

// PeerDetail contains details on peer
type PeerDetail struct {
	ID           string `json:"id"`
	IPAddress    string `json:"ipAddress"`
	MacAddress   string `json:"macAddress"`
	Name         string `json:"name"`
	Organization string `json:"organization"`
}

// Start starts the node
func (node *Node) Start() {
	n := *node

	n.Config = LoadConfig()
	n.ParseIdentity()

	// Load blockchain from storage
	if n.Config.Couch.Enabled == true {
		n.LoadChainFromStorage()
	}

	// TODO
	//n.LoadPair()

	n.DiscoverPeers()

	n.Log(Message{
		Level: 4,
		Text:  "Network started",
	})

	*node = n
}

// DiscoverPeers pulls peer data from downloaded peer lists
func (node *Node) DiscoverPeers() {
	n := *node

	n.DownloadPeerList()

	*node = n
}

// DownloadPeerList pulls peer list from discoveryURL[] config
func (node *Node) DownloadPeerList() {
	n := *node
	list := n.Config.Network.Discovery
	//url := ""

	if len(list) > 0 {
		n.Log(Message{
			Level: 3,
			Text:  "No peers configured for discovery",
		})
		return
	}

	for i := range list {
		n.Log(Message{
			Level: 4,
			Text:  "Downloading peer list from %v" + list[i],
		})
	}
}

// ParseIdentity returns json as struct (TODO!)
func (node *Node) ParseIdentity() {
	n := *node
	c := n.Config.Build

	n.IdentityCertificateLoad()

	file := path.Join(c.IdentityFolder, c.IdentityFile)

	var detail PeerDetail

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

	n.Log(Message{
		Level: 4,
		Text:  "Identity loaded: " + n.Detail.ID,
	})

	*node = n
}
