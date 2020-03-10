package proliferate

import (
	"bytes"
	"encoding/json"

	//	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

//type Object struct {
//	DBName string `json:"db_name"`
//	PurgeSeq string `json:"purge_seq"`
//	UpdateSeq string `json:"update_seq"`
//}

// Record couchDB struct
type Record struct {
	ID           string `json:"_id"`
	Rev          string `json:"_rev"`
	Serial       int    `json:"serial"`
	Timestamp    string `json:"timestamp"`
	Record       string `json:"record"`
	Hash         string `json:"hash"`
	HashPrevious string `json:"hashPrevious"`
}

// RequestBody couchDB struct
type RequestBody struct {
	Method string
	Header http.Header
	Body   string
	Path   string
}

// CouchState couchDB struct
type CouchState struct {
	DBExists bool
}

// CouchDocumentList couchDB struct
type CouchDocumentList struct {
	Rows []CouchDocumentItem `json:"rows"`
}

// CouchDocumentItem couchDB struct
type CouchDocumentItem struct {
	ID  string `json:"id"`
	Key string `json:"key"`
}

// CouchDocumentDetail couchDB struct
type CouchDocumentDetail struct {
	ID        string                 `json:"_id"`
	Serial    int                    `json:"serial"`
	Record    map[string]interface{} `json:"record"`
	Timestamp string                 `json:"timestamp"`
}

// CouchChanges couchDB struct
type CouchChanges struct {
	ID      string              `json:"_id"`
	Results map[int]CouchChange `json:"results"`
}

// CouchChange couchDB struct
type CouchChange struct {
	ID string `json:"id"`
}

// CouchDBDesc couchDB struct
type CouchDBDesc struct {
	DBName             string             `json:"db_name"`
	PurgeSeq           string             `json:"purge_seq"`
	UpdateSeq          string             `json:"update_seq"`
	Sizes              CouchDBDescSizes   `json:"sizes"`
	Other              CouchDBDescOther   `json:"other"`
	DocDeleteCount     int                `json:"doc_del_count"`
	DocCount           int                `json:"doc_count"`
	DiskSize           int                `json:"DiskSize"`
	DiskFormatVersion  int                `json:"disk_format_version"`
	DataSize           int                `json:"data_size"`
	CompactRunning     bool               `json:"compact_running"`
	Cluster            CouchDBDescCluster `json:"cluster"`
	InstantceStartTime string             `json:"instance_start_time"`
}

// CouchDBDescSizes couchDB struct
type CouchDBDescSizes struct {
	File     int `json:"file"`
	External int `json:"external"`
	Active   int `json:"active"`
}

// CouchDBDescOther couchDB struct
type CouchDBDescOther struct {
	DataSize int `json:"data_size"`
}

// CouchDBDescCluster couchDB struct
type CouchDBDescCluster struct {
	Q int `json:"q"`
	N int `json:"n"`
	W int `json:"w"`
	R int `json:"r"`
}

// CouchQueryResults couchDB struct
type CouchQueryResults struct {
	Results []CouchQuerySeq `json:"results"`
}

// CouchQuerySeq couchDB struct
type CouchQuerySeq struct {
	Seq string `json:"seq"`
	ID  string `json:"id"`
}

// CouchQuery couchDB struct
type CouchQuery struct {
	Selector map[string]map[string]string `json:"selector"`
	Fields   []string                     `json:"fields"`
	Sort     map[string]string            `json:"sort"`
	Limit    int                          `json:"limit"`
	Skip     int                          `json:"skip"`
	Stats    bool                         `json:"execution_stats"`
}

/*
// LastBlockFromStorage fetches last block from storage
func (node *Node) LastBlockFromStorage() Block {
	n := *node
	res := n.CouchRaw("_changes?descending=true&limit=1")

	var record CouchQueryResults

	//fmt.Println("----")
	//fmt.Println(res, "\n")
	_ = json.Unmarshal([]byte(res), &record)

	id := record.Results[0].ID
	block := n.LoadBlockFromStorage(id)

	return block
}
*/

// CouchStatus returns couch db description as couchDBDesc
func (node *Node) CouchStatus() CouchDBDesc {
	n := *node
	res := n.CouchRaw("/")
	var status CouchDBDesc

	_ = json.Unmarshal([]byte(res), &status)

	return status
}

// LastBlockFromStorage returns last block from couchDB
func (node *Node) LastBlockFromStorage() Block {
	//n := *node
	var block Block

	//status := n.CouchStatus()

	//serial := status.DocCount - 1

	//jsonQuery := `{"selector":{"serial":` + strconv.Itoa(serial) + `}}`
	selector := `"selector": { "serial": {"$gt":0}}`
	sorter := `"sort": [{"serial": "desc"}]`
	jsonQuery := `{` + selector + `,` + sorter + `}`
	res, _ := node.CouchReq(jsonQuery, "post", "/_find?limit=1")

	fmt.Println("--jsonQuery--")
	fmt.Println(jsonQuery)
	fmt.Println(res)

	json.Unmarshal([]byte(res), &block)

	fmt.Println("--block--")
	fmt.Println(block)

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
	fmt.Println("test", test)

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

// LoadIDsFromStorage returns slice of document IDs TODO at cursor
func (node *Node) LoadIDsFromStorage() []string {
	n := *node
	var docs CouchDocumentList
	//var doc CouchDocumentDetail
	var set []string

	res := n.CouchRaw("/_all_docs")

	json.Unmarshal([]byte(res), &docs)

	for _, v := range docs.Rows {
		set = append(set, v.ID)
	}

	fmt.Println(res)

	return set
}

//LoadChainFromStorage loads blocks from storage
func (node *Node) LoadChainFromStorage() {
	n := *node

	// Limit block load to n.node.memoryLimit
	limit := n.Config.Instance.MemoryRecordLimit + 1
	ls := strconv.Itoa(limit)
	status := n.CouchStatus()

	// Load blocks from last block - memorylimit
	firstSerial := status.DocCount - limit
	arg := strconv.Itoa(firstSerial)

	jsonQuery := `{"selector": { "serial": {"$gt": ` + arg + `}}}`
	res, _ := node.CouchReq(jsonQuery, "post", "/_find?limit="+ls)

	json.Unmarshal([]byte(res), &n)

	*node = n
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

// CouchRaw returns map[string]interface of couch http.Get
func (node *Node) CouchRaw(body string) string {
	n := *node
	request := n.CouchURL() + "/" + body

	//fmt.Println(request)

	response, err := http.Get(request)
	if err != nil {
		n.Log(Message{
			Level: 2,
			Text:  err.Error(),
		})
		return ""
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

// CouchReq returns couch request results as string
func (node *Node) CouchReq(body string, method string, resource string) (string, error) {
	n := *node
	url := n.CouchURL() + resource
	method = strings.ToUpper(method)

	jsonStr := []byte(body)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		//return errors.New(err.Error() + ": " + body)
		fmt.Println(err.Error() + ": " + body)
	}
	defer res.Body.Close()

	responseBody, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(responseBody))

	/*
		if n.StatusCheck(res.Status) == true {
			return errors.New(string(responseBody) + ": " + body)
		}
	*/

	return string(responseBody), nil
}

// StatusCheck couchDB status as bool
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

/*
func (node *Node) CreateDatabase(name string) {
	n := *node

	// TODO exists logic
	n.CouchReq(name, "PUT", "")
}
*/

// CouchPush pushes block to couchDB
func (node *Node) CouchPush(block Block) {
	n := *node
	record := n.MarshalBlock(block)

	//err := n.CouchReq(fmt.Sprintf("%v", block.Record), "POST")
	_, err := n.CouchReq(fmt.Sprintf("%v", record), "POST", "")
	if err != nil {
		n.Log(Message{
			Level: 1,
			Text:  err.Error(),
		})
	}
}

// CouchRequest returns http.Response
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
