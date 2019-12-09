package proliferate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type RequestBody struct {
	Method string
	Header http.Header
	Body   string
}

type CouchState struct {
	DBExists bool
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

// CouchPut http Put body to couchdb
func (node *Node) CouchPut(body string) {
	n := *node
	client := &http.Client{}
	request, err := http.NewRequest("PUT",
		n.CouchURL(),
		strings.NewReader(body))

	//fmt.Println(n.CouchURL(), strings.NewReader(body))
	//fmt.Println(n.CouchURL())
	//fmt.Println("body", body)

	request.ContentLength = int64(len(body))

	response, err := client.Do(request)
	if err != nil {
		n.Log(Message{
			Level: 2,
			Text:  err.Error(),
		})
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)

		if err != nil {
			n.Log(Message{
				Level: 2,
				Text:  err.Error(),
			})
		}
		//fmt.Println("   ", response.StatusCode)

		if response.StatusCode != 200 {
			n.Log(Message{
				Level: 1,
				Text:  string(contents),
			})
		}

		//hdr := response.Header
		//for key, value := range hdr {
		//	fmt.Println("   ", key, ":", value)
		//}
		//fmt.Println(string(contents))
	}
}

// TODO push blocks to couchdb

func (node *Node) CreateDatabase(name string) {
	n := *node

	// TODO exists logic
	n.CouchPut(name)
}

func (node *Node) CouchPush(block Block) {
	n := *node

	req := RequestBody{
		Method: "POST",
		Header: http.Header{},
		Body:   n.ParseRecord(block),
	}

	// TODO Push the block
	fmt.Println(req)
}

func (node *Node) CouchRequest(args RequestBody) *http.Response {
	n := *node
	url := n.CouchURL()
	body := []byte(args.Body)
	client := &http.Client{}

	args.Method = strings.ToUpper(args.Method)

	n.Log(Message{
		Level: 5,
		Text:  "Calling CouchDB @" + url,
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
