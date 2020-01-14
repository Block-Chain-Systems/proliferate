package proliferate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Record struct {
	id           string `json:"_id"`
	rev          string `json:"_rev"`
	serial       int    `json:"serial"`
	timestamp    string `json:"timestamp"`
	record       string `json:"record"`
	hash         string `json:"hash"`
	hashPrevious string `json:"hashPrevious"`
}

type RequestBody struct {
	Method string
	Header http.Header
	Body   string
	Path   string
}

type CouchState struct {
	DBExists bool
}

type CouchDocumentList struct {
	Rows []CouchDocumentDetail `json:"rows"`
}

type CouchDocumentDetail struct {
	ID        string                 `json:"_id"`
	Serial    int                    `json:"serial"`
	Record    map[string]interface{} `json:"record"`
	Timestamp string                 `json:"timestamp"`
}

type CouchChanges struct {
	ID      string              `json:"_id"`
	results map[int]CouchChange `json:"results"`
}

type CouchChange struct {
	ID string `json:"id"`
}

type CouchQueryResults struct {
	Results []CouchQuerySeq `json:"results"`
}

type CouchQuerySeq struct {
	Seq string `json:"seq"`
	ID  string `json:"id"`
}

// LastBlockFromStorage fetches last block from storage
func (node *Node) LastBlockFromStorage() Block {
	n := *node
	res := n.CouchRaw("_changes?descending=true&limit=1")

	var record CouchQueryResults

	fmt.Println("----")
	fmt.Println(res, "\n")
	_ = json.Unmarshal([]byte(res), &record)

	id := record.Results[0].ID
	block := n.LoadBlockFromStorage(id)

	return block
}

// LoadBlockFromStorage returns block from CouchDB by id
func (node *Node) LoadBlockFromStorage(id string) Block {
	n := *node

	var block Block

	res := n.CouchRaw("/" + id)
	_ = json.Unmarshal([]byte(res), &block)

	return block
}

// CouchURL parses config.json and returns couch http URL
func (node *Node) CouchURL() string {
	n := *node
	c := n.Config.Couch
	return c.Protocol + "://" + c.Host + ":" + c.Port + "/" + c.Database
}

// CouchTest Verifies CouchDB host is responsive
func (node *Node) CouchTest() bool {
	n := *node
	data := n.CouchGet("")

	if data["couchdb"] == "Welcome" {
		return true
	}

	return false
}

// DBExists returns true if couchDB is found
func (node *Node) DBExists() bool {
	n := *node
	c := n.Config.Couch

	if c.State.DBExists == true {
		return true
	}

	res := n.CouchRequest(RequestBody{
		Method: "GET",
		Body:   "/" + c.Database,
	})

	if res == nil {
		return false
	}

	if res.StatusCode == 200 {
		n.Config.Couch.State.DBExists = true
		return true
	}

	defer res.Body.Close()
	test, _ := ioutil.ReadAll(res.Body)

	fmt.Println(test)

	//for _, v := range list {
	//	fmt.Println(v)
	//}

	*node = n
	return false
}

//func (node *Node) InitialzeDatabase() {
//	n := *node
//	n.CouchPut("proliferate")
//}

// CouchGet returns map[string]interface of couch http.Get
func (node *Node) CouchGet(body string) map[string]interface{} {
	n := *node
	request := n.CouchURL() + "/" + body

	fmt.Println(request)

	response, err := http.Get(request)
	if err != nil {
		n.Log(Message{
			Level: 2,
			Text:  err.Error(),
		})
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		n.Log(Message{
			Level: 2,
			Text:  err.Error(),
		})
	}

	data := make(map[string]interface{})

	_ = json.Unmarshal([]byte(responseData), &data)

	return data
}

// Returns slice of document IDs TODO at cursor
func (node *Node) LoadIDsFromStorage() []string {
	n := *node
	var docs CouchDocumentList
	var set []string

	res := n.CouchRaw("/_all_docs")
	json.Unmarshal([]byte(res), &docs)

	for _, v := range docs.Rows {
		set = append(set, v.ID)
	}

	return set
}

func (node *Node) LoadChainFromStorage() {
	n := *node

	ids := n.LoadIDsFromStorage()

	for _, v := range ids {
		//fmt.Println(i, v)
		fmt.Println(n.CouchRaw("/" + v))
		// TODO Structure the records you are here!
	}
}

/*
//  LoadDocumentsFromStorage returns documents TODO at cursor
func (node *Node) LoadDocumentsFromStorage() CouchDocumentList {
	n := *node
	var docs CouchDocumentList

	if n.DBExists() != true {
		n.Log(Message{
			Level: 1,
			Text:  "Unable to load chain from storage",
		})

		return docs
	}

	res := n.CouchRaw("/_all_docs?include_docs=true")
	json.Unmarshal([]byte(res), &docs)

	return docs
}
*/

// CouchGet returns map[string]interface of couch http.Get
func (node *Node) CouchRaw(body string) string {
	n := *node
	request := n.CouchURL() + "/" + body

	response, err := http.Get(request)
	if err != nil {
		n.Log(Message{
			Level: 2,
			Text:  err.Error(),
		})
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		n.Log(Message{
			Level: 2,
			Text:  err.Error(),
		})
	}

	//data := make(map[string]interface{})
	//_ = json.Unmarshal([]byte(responseData), &data)
	return string(responseData)
}

func (node *Node) CouchReq(body string, method string) error {
	n := *node
	url := n.CouchURL()

	jsonStr := []byte(body)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return errors.New(err.Error() + ": " + body)
	}
	defer res.Body.Close()

	responseBody, _ := ioutil.ReadAll(res.Body)

	if n.StatusCheck(res.Status) == true {
		return errors.New(string(responseBody) + ": " + body)
	}

	return nil
}

func (node *Node) StatusCheck(code string) bool {
	n := *node
	code = code[0:3]

	value, err := strconv.Atoi(code)
	if err != nil {
		n.Log(Message{
			Level: 2,
			Text:  err.Error(),
		})
	}

	if value >= 200 && value <= 202 {
		return true
	}

	return false
}

func (node *Node) CreateDatabase(name string) {
	n := *node

	// TODO exists logic
	n.CouchReq(name, "PUT")
}

func (node *Node) CouchPush(block Block) {
	n := *node
	record := n.MarshalBlock(block)

	//err := n.CouchReq(fmt.Sprintf("%v", block.Record), "POST")
	err := n.CouchReq(fmt.Sprintf("%v", record), "POST")
	if err != nil {
		n.Log(Message{
			Level: 1,
			Text:  err.Error(),
		})
	}
}

func (node *Node) CouchRequest(args RequestBody) *http.Response {
	n := *node
	url := n.CouchURL()
	body := []byte(args.Body)
	client := &http.Client{}

	args.Method = strings.ToUpper(args.Method)

	n.Log(Message{
		Level: 5,
		Text:  "Calling CouchDB: " + url,
	})

	req, reqErr := http.NewRequest(args.Method, url, bytes.NewBuffer(body))

	if reqErr != nil {
		n.Log(Message{
			Level: 2,
			Text:  reqErr.Error(),
		})
	}

	res, resErr := client.Do(req)

	if resErr != nil {
		n.Log(Message{
			Level: 2,
			Text:  resErr.Error(),
		})
	}

	//fmt.Println(res)
	return res
}
